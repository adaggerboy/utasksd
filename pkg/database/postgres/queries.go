package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/adaggerboy/utasksd/models/generic"
	"github.com/adaggerboy/utasksd/models/reqresp"
)

//go:embed schema.sql
var schemaCreateQuery string

func getPreparedMap() map[string]string {
	return map[string]string{
		"create_task": `insert into tasks (assigner_id, project_id, name, status, priority, start_date, due_date, description) values ($1, $2, $3, $4, $5, $6, $7, $8) on conflict do nothing returning id`,
		"read_task":   `select assigner_id, project_id, name, status, priority, start_date, due_date, description from tasks where id = $1`,
		"update_task": `update tasks set 
			assigner_id = coalesce($2, assigner_id),
			project_id = coalesce($3, project_id),
			name = coalesce($4, name),
			status = coalesce($5, status),
			priority = coalesce($6, priority),
			start_date = coalesce($7, start_date),
			due_date = coalesce($8, due_date),
			description = coalesce($9, description)
			where id = $1
		`,
		"delete_task": `delete from tasks where id = $1`,

		"delete_task_assignees": `delete from tasks_assignees where task_id = $1;`,
		"delete_task_assignee":  `delete from tasks_assignees where task_id = $1 and assignee_id = $2;`,
		"add_task_assignee":     `insert into tasks_assignees (task_id, assignee_id) values ($1, $2) on conflict do nothing`,
		"get_task_assignees":    `select assignee_id from tasks_assignees where task_id = $1`,

		"delete_task_attachments": `delete from tasks_attachments where task_id = $1`,
		"delete_task_attachment":  `delete from tasks_attachments where task_id = $1 and attachment = $2 `,
		"add_task_attachment":     `insert into tasks_attachments (task_id, attachment) values ($1, $2) on conflict do nothing`,
		"get_task_attachments":    `select attachment from tasks_attachments where task_id = $1`,

		"delete_task_issues": `delete from tasks_issues where task_id = $1`,
		"delete_task_issue":  `delete from tasks_issues where task_id = $1 and issue_id = $2 `,
		"add_task_issue":     `insert into tasks_issues (task_id, issue_id) values ($1, $2) on conflict do nothing`,
		"get_task_issues":    `select issue_id from tasks_issues where task_id = $1`,

		"delete_dependent_tasks": `delete from tasks_dependencies where task_id = $1`,
		"delete_dependent_task":  `delete from tasks_dependencies where task_id = $1 and dep_task_id = $2 `,
		"add_dependent_task":     `insert into tasks_dependencies (task_id, dep_task_id, dependency_type) values($1, $2, $3) on conflict do nothing`,
		"get_dependent_tasks":    `select dep_task_id, dependency_type from tasks_dependencies where task_id = $1`,

		"create_comment": `insert into comments (task_id, user_id, text) values ($1, $2, $3) returning id`,
		"read_comment":   `select task_id, user_id, text, date from comments where id = $1`,
		"delete_comment": `delete from comments where id = $1`,

		"get_task_comments": `select id, user_id, text, date from comments where task_id = $1`,

		"add_tracking_record":  "insert into tracking_records (task_id, text, seconds) values ($1, $2, $3)",
		"get_tracking_records": "select text, seconds, end_date from tracking_records where task_id = $1",

		"create_project": `insert into projects (owner_id, name, logo, description) values ($1, $2, $3, $4) on conflict do nothing returning id`,
		"read_project":   `select owner_id, name, logo, description from projects where id = $1`,
		"update_project": `update projects set 
			owner_id = coalesce($2, owner_id),
			name = coalesce($3, name),
			logo = coalesce($4, logo),
			description = coalesce($5, description)
			where id = $1
			`,
		"delete_project": `delete from projects where id = $1`,

		"get_project_members":      `select user_id, access_level from users_projects where project_id = $1`,
		"get_project_member_level": `select access_level from users_projects where project_id = $1 and user_id = $2`,
		"delete_project_members":   `delete from users_projects where project_id = $1 and access_level = $2`,
		"add_project_member":       `insert into users_projects (project_id, access_level, user_id) values ($1, $2, $3)`,

		"search_tasks": `select distinct id from tasks 
			inner join users_projects on tasks.project_id = users_projects.project_id
			left join tasks_assignees on tasks_assignees.task_id = tasks.id 
			left join tasks_issues on tasks_issues.task_id = tasks.id where 
			tasks.project_id = coalesce($1, tasks.project_id) and
			assigner_id = coalesce($2, assigner_id) and
			(assignee_id = coalesce($3, assignee_id) or 
				assignee_id is null and $3 is null) and
			user_id = $12 and
			(issue_id = coalesce($4, issue_id) or 
				issue_id is null and $4 is null) and
			status = coalesce($5, status) and
			priority = coalesce($6, priority) and (
				$7 = false or
				name like '%' || coalesce($8, '') || '%' or
				description like '%' || coalesce($8, '') || '%'
			) and (
				$9 = false or
				(start_date >= coalesce($10, start_date) and
				start_date <= coalesce($11, start_date) and (
					due_date is null and $11 is null or 
					due_date is null and $11 > now() or
					due_date >= coalesce($10, due_date) and
					due_date <= coalesce($11, due_date)
				))
			)
		`,

		"create_user": `insert into users (username, email, firstname, lastname, avatar) values ($1, $2, $3, $4, $5) on conflict do nothing returning id`,
		"read_user":   `select username, email, firstname, lastname, avatar, is_active, is_admin, is_director from users where id = $1`,
		"update_user": `update users set 
			username = coalesce($2, username),
			email = coalesce($3, email),
			firstname = coalesce($4, firstname),
			lastname = coalesce($5, lastname),
			avatar = coalesce($6, avatar)
			where id = $1
		`,
		"set_user_permissions": `update users set 
			is_active = $2,
			is_admin = $3,
			is_director = $4
			where id = $1
		`,
		"get_all_users": `select id, username, email, firstname, lastname, avatar, is_active, is_admin, is_director from users`,

		"change_activation_status": `update users set is_active = $2 where id = $1`,
		"get_user_by_username":     `select id, is_active from users where username = $1`,

		"get_project_tasks": `select id from tasks where project_id = $1`,
		"get_user_projects": `select project_id from users_projects where user_id = $1`,
		"get_user_project_tasks": `select task_id from tasks_assignees inner join tasks on task_id = tasks.id where assignee_id = $1 and project_id = $2
			union
			select id from tasks where assigner_id = $1 and project_id = $2
		`,
		"get_user_tasks": `select task_id from tasks_assignees where assignee_id = $1
			union
			select id from tasks where assigner_id = $1
		`,

		"get_project_issues": `select id from issues where project_id = $1`,
		"get_user_issues":    `select issues.id from issues inner join users_projects on users_projects.project_id = issues.project_id where users_projects.user_id = $1`,

		"get_task_user_capabilities": `select 'assigner' as role from tasks where id = $1 and assigner_id = $2
			union
			select 'assignee' as role from tasks_assignees where task_id = $1 and assignee_id = $2
			union
			select 'watcher' as role from users_projects inner join tasks on tasks.project_id = users_projects.project_id where tasks.id = $1 and user_id = $2
			union
			select 'owner' as role from projects inner join tasks on tasks.project_id = projects.id where tasks.id = $1 and owner_id = $2
		`,

		"get_issue_user_capabilities": `select 'editor' as role from issues where id = $1 and publisher_id = $2
			union
			select 'watcher' as role from users_projects inner join issues on issues.project_id = users_projects.project_id where issues.id = $1 and user_id = $2
			union
			select 'owner' as role from projects inner join issues on issues.project_id = projects.id where issues.id = $1 and owner_id = $2
		`,

		"get_user_credentials":    `select pass_hash, salt, algo from user_credentials where id = $1`,
		"update_user_credentials": `insert into user_credentials (id, pass_hash, salt, algo) values ($1, $2, $3, $4) on conflict (id) do update set pass_hash = $2, salt = $3, algo = $4`,

		"insert_attachment": "insert into attachments (cdn_path) values ($1)",

		"create_issue": `insert into issues (name, publisher_id, project_id, description, reporter) values ($1, $2, $3, $4, $5) on conflict do nothing returning id`,
		"read_issue":   `select publisher_id, project_id, name, status, reporter, open_date, close_date, description from issues where id = $1`,
		"update_issue": `update issues set 
			name = coalesce($2, name),
			project_id = coalesce($3, project_id),
			description = coalesce($4, description),
			status = coalesce($5, status)
			where id = $1
		`,
		"delete_issue": `delete from issues where id = $1`,

		"search_issues": `select distinct id from issues 
			inner join users_projects on issues.project_id = users_projects.project_id where
			issues.project_id = coalesce($1, issues.project_id) and
			publisher_id = coalesce($2, publisher_id) and
			user_id = $9 and
			status = coalesce($3, status) and (
				$4 = false or
				name like '%' || coalesce($5, '') || '%' or
				description like '%' || coalesce($5, '') || '%' or
				reporter like '%' || coalesce($5, '') || '%'
			) and (
				$6 = false or
				(open_date >= coalesce($7, open_date) and
				open_date <= coalesce($8, open_date) and (
					close_date is null and $8 is null or 
					close_date is null and $8 > now() or
					close_date >= coalesce($7, close_date) and
					close_date <= coalesce($8, close_date)
				))
			)
		`,
		"delete_issue_attachments": `delete from issues_attachments where issue_id = $1`,
		"delete_issue_attachment":  `delete from issues_attachments where issue_id = $1 and attachment = $2 `,
		"add_issue_attachment":     `insert into issues_attachments (issue_id, attachment) values ($1, $2) on conflict do nothing`,
		"get_issue_attachments":    `select attachment from issues_attachments where issue_id = $1`,

		"tasks_report":           `select * from tasks_report($1, $2, $3)`,
		"time_efficiency_report": `select * from time_efficiency_report($1, $2, $3)`,
	}
}

func (db *PostgresConnection) DeploySchema() error {
	_, err := db.db.Exec(schemaCreateQuery)
	return err
}

func (db *PostgresConnection) prepareStmt(ctx context.Context, queryNames []string) (stmts map[string]*sql.Stmt, err error) {
	stmts = map[string]*sql.Stmt{}
	for _, v := range queryNames {
		var stmt *sql.Stmt
		stmt, err = db.tx.PrepareContext(ctx, getPreparedMap()[v])
		if err != nil {
			return
		}
		stmts[v] = stmt
	}
	return
}

func closeStmt(stmts map[string]*sql.Stmt, _ error) (err error) {
	for _, v := range stmts {
		err = v.Close()
		if err != nil {
			return
		}
	}
	return
}

func (db *PostgresConnection) genericExec(ctx context.Context, stmtName string, args ...any) (err error) {
	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)

	_, err = stmts[stmtName].ExecContext(ctx, args...)
	return
}

func (db *PostgresConnection) genericQuery(ctx context.Context, stmtName string, args []any, rets ...any) (err error) {
	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)

	rows, err := stmts[stmtName].QueryContext(ctx, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(rets...)
	}
	return
}

func (db *PostgresConnection) genericQueryForeach(ctx context.Context, stmtName string, foreach func(*sql.Rows) error, args ...any) (err error) {
	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)

	rows, err := stmts[stmtName].QueryContext(ctx, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = foreach(rows)
		if err != nil {
			return
		}
	}
	return
}

func genericExecForeach[V any](db *PostgresConnection, ctx context.Context, stmtName string, instances []V, arg any) (err error) {
	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)
	for _, v := range instances {
		val := V(v)
		_, err = stmts[stmtName].ExecContext(ctx, arg, val)
		if err != nil {
			return
		}
	}
	return
}

func genericExecForeach2[V any](db *PostgresConnection, ctx context.Context, stmtName string, instances []V, arg1 any, arg2 any) (err error) {
	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)
	for _, v := range instances {
		val := V(v)
		_, err = stmts[stmtName].ExecContext(ctx, arg1, arg2, val)
		if err != nil {
			return
		}
	}
	return
}

func (db *PostgresConnection) genericExecReturnID(ctx context.Context, stmtName string, args ...any) (id *int, err error) {
	id = nil

	stmts, err := db.prepareStmt(ctx, []string{
		stmtName,
	})
	if err != nil {
		return
	}
	defer closeStmt(stmts, err)

	rows, err := stmts[stmtName].QueryContext(ctx, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		id = new(int)
		err = rows.Scan(id)
	}
	return
}

func (db *PostgresConnection) CreateTask(ctx context.Context, task generic.Task) (taskID *int, err error) {
	if taskID, err = db.genericExecReturnID(ctx, "create_task", task.Assigner, task.Project, task.Name, task.Status, task.Priority, task.StartDate, task.DueDate, task.Description); err != nil {
		return
	}

	if task.Assignees != nil {
		if err = genericExecForeach(db, ctx, "add_task_assignee", *task.Assignees, taskID); err != nil {
			return
		}
	}
	if task.Attachments != nil {
		if err = genericExecForeach(db, ctx, "add_task_attachment", *task.Attachments, taskID); err != nil {
			return
		}
	}
	if task.LinkedIssues != nil {
		if err = genericExecForeach(db, ctx, "add_task_issue", *task.LinkedIssues, taskID); err != nil {
			return
		}
	}
	if task.DependentTasks != nil {
		for _, v := range *task.DependentTasks {
			if err = db.genericExec(ctx, "add_dependent_task", taskID, v.ID, v.TaskDependencyType); err != nil {
				return
			}
		}
	}
	return
}

func (db *PostgresConnection) UpdateTask(ctx context.Context, taskID int, task generic.Task) (err error) {

	if err = db.genericExec(ctx, "update_task", taskID, task.Assigner, task.Project, task.Name, task.Status, task.Priority, task.StartDate, task.DueDate, task.Description); err != nil {
		return
	}
	if task.Assignees != nil {
		if err = db.genericExec(ctx, "delete_task_assignees", taskID); err != nil {
			return
		}
		if err = genericExecForeach(db, ctx, "add_task_assignee", *task.Assignees, taskID); err != nil {
			return
		}
	}
	if task.Attachments != nil {
		if err = db.genericExec(ctx, "delete_task_attachments", taskID); err != nil {
			return
		}
		if err = genericExecForeach(db, ctx, "add_task_attachment", *task.Attachments, taskID); err != nil {
			return
		}
	}
	if task.DependentTasks != nil {
		if err = db.genericExec(ctx, "delete_dependent_tasks", taskID); err != nil {
			return
		}
		for _, v := range *task.DependentTasks {
			if err = db.genericExec(ctx, "add_dependent_task", taskID, v.ID, v.TaskDependencyType); err != nil {
				return
			}
		}
	}
	if task.LinkedIssues != nil {
		if err = db.genericExec(ctx, "delete_task_issues", taskID); err != nil {
			return
		}
		if err = genericExecForeach(db, ctx, "add_task_issue", *task.LinkedIssues, taskID); err != nil {
			return
		}
	}
	return
}

func (db *PostgresConnection) DeleteTask(ctx context.Context, taskID int) (err error) {
	return db.genericExec(ctx, "delete_task", taskID)
}

func (db *PostgresConnection) ReadTask(ctx context.Context, taskID int) (task *generic.Task, err error) {
	task = generic.AllocateTask()
	err = db.genericQuery(ctx, "read_task", []any{taskID}, task.Assigner, task.Project, task.Name, task.Status, task.Priority, task.StartDate, task.DueDate, task.Description)
	return
}
func (db *PostgresConnection) ReadTaskRecursive(ctx context.Context, taskID int) (task *generic.Task, err error) {
	task, err = db.ReadTask(ctx, taskID)
	if err != nil || task == nil {
		return
	}
	if err = db.genericQueryForeach(ctx, "get_task_issues", func(r *sql.Rows) error {
		var item int
		if err := r.Scan(&item); err != nil {
			return err
		}
		(*task.LinkedIssues) = append(*task.LinkedIssues, item)
		return nil
	}, taskID); err != nil {
		return
	}

	if err = db.genericQueryForeach(ctx, "get_task_comments", func(r *sql.Rows) error {
		comment := generic.AllocateComment()
		var date time.Time
		if err := r.Scan(comment.ID, comment.AuthorID, comment.Text, &date); err != nil {
			return err
		}
		timeAgo := time.Since(date)
		comment.TimeAgo = &timeAgo
		(*task.Comments) = append(*task.Comments, *comment)
		return nil
	}, taskID); err != nil {
		return
	}

	if err = db.genericQueryForeach(ctx, "get_tracking_records", func(r *sql.Rows) error {
		tr := generic.AllocateTrackingRecord()
		var seconds int
		if err := r.Scan(tr.Text, &seconds, tr.EndDate); err != nil {
			return err
		}
		tr.Duration = &generic.CommonDuration{Duration: time.Second * time.Duration(seconds)}
		(*task.TrackedRecords) = append(*task.TrackedRecords, *tr)
		return nil
	}, taskID); err != nil {
		return
	}

	if err = db.genericQueryForeach(ctx, "get_task_assignees", func(r *sql.Rows) error {
		var assignee int
		if err := r.Scan(&assignee); err != nil {
			return err
		}
		(*task.Assignees) = append(*task.Assignees, assignee)
		return nil
	}, taskID); err != nil {
		return
	}

	if err = db.genericQueryForeach(ctx, "get_task_attachments", func(r *sql.Rows) error {
		var att string
		if err := r.Scan(&att); err != nil {
			return err
		}
		(*task.Attachments) = append(*task.Attachments, att)
		return nil
	}, taskID); err != nil {
		return
	}

	err = db.genericQueryForeach(ctx, "get_dependent_tasks", func(r *sql.Rows) error {
		var dtask generic.DependentTask
		if err := r.Scan(&dtask.ID, &dtask.TaskDependencyType); err != nil {
			return err
		}
		(*task.DependentTasks) = append(*task.DependentTasks, dtask)
		return nil
	}, taskID)

	return
}

func (db *PostgresConnection) CreateUser(ctx context.Context, user generic.User) (userID *int, err error) {
	return db.genericExecReturnID(ctx, "create_user", user.Username, user.Email, user.Firstname, user.Lastname, user.AvatarPath)
}
func (db *PostgresConnection) GetUserByUsername(ctx context.Context, username string) (userID *int, isActive bool, err error) {
	userID = new(int)
	err = db.genericQuery(ctx, "get_user_by_username", []any{username}, userID, &isActive)
	return
}
func (db *PostgresConnection) UpdateUser(ctx context.Context, userID int, user generic.User) (err error) {
	return db.genericExec(ctx, "update_user", userID, user.Username, user.Email, user.Firstname, user.Lastname, user.AvatarPath)
}
func (db *PostgresConnection) DeactivateUser(ctx context.Context, userID int) (err error) {
	return db.genericExec(ctx, "change_activation_status", userID, false)
}

func (db *PostgresConnection) ReadUser(ctx context.Context, userID int) (user *generic.User, err error) {
	user = generic.AllocateUser()
	err = db.genericQuery(ctx, "read_user", []any{userID}, user.Username, user.Email, user.Firstname, user.Lastname, user.AvatarPath, user.IsActive, user.IsAdmin, user.IsDirector)
	return
}

func (db *PostgresConnection) CreateProject(ctx context.Context, project generic.Project) (projectID *int, err error) {

	if projectID, err = db.genericExecReturnID(ctx, "create_project", project.OwnerID, project.Name, project.LogoPath, project.Description); err != nil {
		return
	}
	if project.Workers != nil {
		if err = genericExecForeach2(db, ctx, "add_project_member", *project.Workers, projectID, "worker"); err != nil {
			return
		}
	}
	if project.Managers != nil {
		if err = genericExecForeach2(db, ctx, "add_project_member", *project.Managers, projectID, "manager"); err != nil {
			return
		}
	}
	if project.SupportAgents != nil {
		err = genericExecForeach2(db, ctx, "add_project_member", *project.SupportAgents, projectID, "support")
	}
	return
}
func (db *PostgresConnection) UpdateProject(ctx context.Context, projectID int, project generic.Project) (err error) {
	if err = db.genericExec(ctx, "update_project", projectID, project.OwnerID, project.Name, project.LogoPath, project.Description); err != nil {
		return
	}
	if project.Workers != nil {
		if err = db.genericExec(ctx, "delete_project_members", projectID, "worker"); err != nil {
			return
		}
	}
	if project.Managers != nil {
		if err = db.genericExec(ctx, "delete_project_members", projectID, "manager"); err != nil {
			return
		}
	}
	if project.SupportAgents != nil {
		if err = db.genericExec(ctx, "delete_project_members", projectID, "support"); err != nil {
			return
		}
	}
	if project.Workers != nil {
		if err = genericExecForeach2(db, ctx, "add_project_member", *project.Workers, projectID, "worker"); err != nil {
			return
		}
	}
	if project.Managers != nil {
		if err = genericExecForeach2(db, ctx, "add_project_member", *project.Managers, projectID, "manager"); err != nil {
			return
		}
	}
	if project.SupportAgents != nil {
		if err = genericExecForeach2(db, ctx, "add_project_member", *project.SupportAgents, projectID, "support"); err != nil {
			return
		}
	}
	return
}
func (db *PostgresConnection) DeleteProject(ctx context.Context, projectID int) (err error) {
	return db.genericExec(ctx, "delete_project", projectID)
}
func (db *PostgresConnection) ReadProject(ctx context.Context, projectID int) (project *generic.Project, err error) {

	project = generic.AllocateProject()
	err = db.genericQuery(ctx, "read_project", []any{projectID}, project.OwnerID, project.Name, project.LogoPath, project.Description)
	return
}
func (db *PostgresConnection) ReadProjectRecursive(ctx context.Context, projectID int) (project *generic.Project, err error) {
	project, err = db.ReadProject(ctx, projectID)
	if err != nil || project == nil {
		return
	}
	return project, db.genericQueryForeach(ctx, "get_project_members", func(r *sql.Rows) error {
		var id int
		var level string
		if err := r.Scan(&id, &level); err != nil {
			return err
		}
		switch level {
		case "worker":
			(*project.Workers) = append((*project.Workers), id)
		case "support":
			(*project.SupportAgents) = append((*project.SupportAgents), id)
		case "manager":
			(*project.Managers) = append((*project.Managers), id)
		case "owner":
		default:
			return fmt.Errorf("unexpected user level: %s", level)
		}
		return nil
	}, projectID)
}

func (db *PostgresConnection) EnsureTaskCapability(ctx context.Context, userID, taskID int) (readCapability, editCapability, manageCapability bool, err error) {
	readCapability = false
	editCapability = false
	manageCapability = false

	err = db.genericQueryForeach(ctx, "get_task_user_capabilities", func(r *sql.Rows) error {
		var level string
		if err := r.Scan(&level); err != nil {
			return err
		}
		switch level {
		case "assigner":
			readCapability = true
			editCapability = true
			manageCapability = true
		case "assignee":
			readCapability = true
			editCapability = true
		case "watcher":
			readCapability = true
		case "owner":
			editCapability = true
			readCapability = true
			manageCapability = true
		default:
			return fmt.Errorf("unknown task access level: %s", level)
		}
		return nil
	}, taskID, userID)
	return
}

func (db *PostgresConnection) EnsureIssueCapability(ctx context.Context, userID, issueID int) (readCapability, editCapability bool, err error) {
	readCapability = false
	editCapability = false

	err = db.genericQueryForeach(ctx, "get_issue_user_capabilities", func(r *sql.Rows) error {
		var level string
		if err := r.Scan(&level); err != nil {
			return err
		}
		switch level {
		case "assigner":
		case "editor":
			readCapability = true
			editCapability = true
		case "watcher":
			readCapability = true
		case "owner":
			editCapability = true
			readCapability = true
		default:
			return fmt.Errorf("unknown task access level: %s", level)
		}
		return nil
	}, issueID, userID)
	return
}

func (db *PostgresConnection) EnsureUsersCapability(ctx context.Context, userID int) (manageCapability, directorCapability bool, err error) {
	manageCapability = false
	user, err := db.ReadUser(ctx, userID)
	if err != nil || user == nil {
		return
	}
	manageCapability = *user.IsAdmin && *user.IsActive
	directorCapability = (*user.IsAdmin || *user.IsDirector) && *user.IsActive
	return
}
func (db *PostgresConnection) EnsureProjectCapability(ctx context.Context, userID, projectID int) (readCapability, tasksCapability, issuesCapability, manageCapability bool, err error) {
	readCapability = false
	tasksCapability = false
	issuesCapability = false
	manageCapability = false
	var level string

	if err = db.genericQuery(ctx, "get_project_member_level", []any{projectID, userID}, &level); err != nil {
		return
	}

	switch level {
	case "worker":
		readCapability = true
	case "support":
		readCapability = true
		issuesCapability = true
	case "manager":
		readCapability = true
		tasksCapability = true
	case "owner":
		readCapability = true
		tasksCapability = true
		issuesCapability = true
		manageCapability = true
	default:
		err = fmt.Errorf("unexpected user level: %s", level)
		return
	}
	return
}

func (db *PostgresConnection) GetUserProjects(ctx context.Context, userID int) (projectIDs []int, err error) {
	projectIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_user_projects", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		projectIDs = append(projectIDs, next)
		return nil
	}, userID)
	return
}
func (db *PostgresConnection) GetUserTasks(ctx context.Context, userID int) (taskIDs []int, err error) {
	taskIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_user_tasks", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		taskIDs = append(taskIDs, next)
		return nil
	}, userID)
	return
}
func (db *PostgresConnection) GetUserProjectTasks(ctx context.Context, userID, projectID int) (taskIDs []int, err error) {
	taskIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_user_project_tasks", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		taskIDs = append(taskIDs, next)
		return nil
	}, userID, projectID)
	return
}
func (db *PostgresConnection) GetProjectTasks(ctx context.Context, projectID int) (taskIDs []int, err error) {
	taskIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_project_tasks", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		taskIDs = append(taskIDs, next)
		return nil
	}, projectID)
	return
}

func (db *PostgresConnection) SearchTasks(ctx context.Context, query reqresp.SearchTasks, userID int) (taskIDs []int, err error) {
	taskIDs = []int{}
	byName := query.Name != nil
	byDate := query.DueDate != nil || query.StartDate != nil
	err = db.genericQueryForeach(ctx, "search_tasks", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		taskIDs = append(taskIDs, next)
		return nil
	}, query.Project, query.Assigner, query.Assignee, query.Issue, query.Status, query.Priority, byName, query.Name, byDate, query.StartDate, query.DueDate, userID)
	return
}

func (db *PostgresConnection) RegisterAttachment(ctx context.Context, attachment string) (err error) {
	return db.genericExec(ctx, "insert_attachment", attachment)
}

func (db *PostgresConnection) ReadComment(ctx context.Context, commentID int) (comment *generic.Comment, err error) {
	comment = generic.AllocateComment()
	var date time.Time
	if err = db.genericQuery(ctx, "read_comment", []any{commentID}, comment.TaskID, comment.AuthorID, comment.Text, &date); err != nil {
		return
	}
	timeAgo := time.Since(date)
	comment.TimeAgo = &timeAgo
	return
}
func (db *PostgresConnection) PublishComment(ctx context.Context, userID int, taskID int, text string) (commentID *int, err error) {
	return db.genericExecReturnID(ctx, "create_comment", taskID, userID, text)
}
func (db *PostgresConnection) DeleteComment(ctx context.Context, commentID int) (err error) {
	return db.genericExec(ctx, "delete_comment", commentID)
}

func (db *PostgresConnection) TrackTaskActivity(ctx context.Context, taskID int, text string, duration int) (err error) {
	return db.genericExec(ctx, "add_tracking_record", taskID, text, duration)
}

func (db *PostgresConnection) CreateIssue(ctx context.Context, issue generic.Issue) (issueID *int, err error) {
	if issueID, err = db.genericExecReturnID(ctx, "create_issue", issue.Name, issue.Registrar, issue.Project, issue.Description, issue.Reporter); err != nil {
		return
	}
	if issue.Attachments != nil {
		err = genericExecForeach(db, ctx, "add_issue_attachment", *issue.Attachments, issueID)
	}
	return
}
func (db *PostgresConnection) UpdateIssue(ctx context.Context, issueID int, issue generic.Issue) (err error) {
	if err = db.genericExec(ctx, "update_issue", issueID, issue.Name, issue.Project, issue.Description, issue.Status); err != nil {
		return
	}
	if issue.Attachments != nil {
		if err = db.genericExec(ctx, "delete_issue_attachments", issueID); err != nil {
			return
		}
		err = genericExecForeach(db, ctx, "add_issue_attachment", *issue.Attachments, issueID)
	}
	return
}
func (db *PostgresConnection) DeleteIssue(ctx context.Context, issueID int) (err error) {
	return db.genericExec(ctx, "delete_issue", issueID)

}
func (db *PostgresConnection) ReadIssue(ctx context.Context, issueID int) (issue *generic.Issue, err error) {
	issue = generic.AllocateIssue()
	return issue, db.genericQuery(ctx, "read_issue", []any{issueID}, issue.Registrar, issue.Project, issue.Name, issue.Status, &issue.Reporter, issue.StartDate, &issue.CloseDate, issue.Description)
}

func (db *PostgresConnection) ReadIssueRecursive(ctx context.Context, issueID int) (issue *generic.Issue, err error) {
	issue, err = db.ReadIssue(ctx, issueID)
	if err != nil || issue == nil {
		return
	}

	return issue, db.genericQueryForeach(ctx, "get_issue_attachments", func(r *sql.Rows) error {
		var att string
		if err := r.Scan(&att); err != nil {
			return err
		}
		(*issue.Attachments) = append(*issue.Attachments, att)
		return nil
	}, issueID)
}

func (db *PostgresConnection) SearchIssues(ctx context.Context, query reqresp.SearchIssues, userID int) (issueIDs []int, err error) {
	issueIDs = []int{}
	byName := query.Name != nil
	byDate := query.DueDate != nil || query.StartDate != nil
	err = db.genericQueryForeach(ctx, "search_issues", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		issueIDs = append(issueIDs, next)
		return nil
	}, query.Project, query.Registrar, query.Status, byName, query.Name, byDate, query.StartDate, query.DueDate, userID)
	return
}

func (db *PostgresConnection) GetUserIssues(ctx context.Context, userID int) (issueIDs []int, err error) {
	issueIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_user_issues", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		issueIDs = append(issueIDs, next)
		return nil
	}, userID)
	return
}
func (db *PostgresConnection) GetProjectIssues(ctx context.Context, projectID int) (issueIDs []int, err error) {
	issueIDs = []int{}
	err = db.genericQueryForeach(ctx, "get_project_issues", func(r *sql.Rows) error {
		var next int
		if err := r.Scan(&next); err != nil {
			return err
		}
		issueIDs = append(issueIDs, next)
		return nil
	}, projectID)
	return
}

func (db *PostgresConnection) Close() (err error) {
	err = db.tx.Commit()
	if err != nil {
		return
	}
	return db.db.Close()
}

func (db *PostgresConnection) RollbackClose() (err error) {
	err = db.tx.Rollback()
	if err != nil {
		return
	}
	return db.db.Close()
}

func (db *PostgresConnection) GetAllUsers(ctx context.Context) (users []generic.User, err error) {
	users = []generic.User{}
	err = db.genericQueryForeach(ctx, "get_all_users", func(r *sql.Rows) error {
		user := generic.AllocateUser()
		if err := r.Scan(&user.ID, user.Username, user.Email, user.Firstname, user.Lastname, user.AvatarPath, user.IsActive, user.IsAdmin, user.IsDirector); err != nil {
			return err
		}
		users = append(users, *user)
		return nil
	})
	return
}

func (db *PostgresConnection) SetUserPermissions(ctx context.Context, userID int, username string, isActive, isAdmin, isDirector bool) (err error) {
	if err = db.SetRole(ctx, username, "administrator", isAdmin); err != nil {
		return
	}
	if err = db.SetRole(ctx, username, "director", isDirector); err != nil {
		return
	}
	return db.genericExec(ctx, "set_user_permissions", userID, isActive, isAdmin, isDirector)
}

func (db *PostgresConnection) TasksReport(ctx context.Context, request reqresp.ReportRequest) (perUser []generic.TasksReport, err error) {
	perUser = []generic.TasksReport{}
	err = db.genericQueryForeach(ctx, "tasks_report", func(r *sql.Rows) error {
		var next generic.TasksReport
		if err := r.Scan(&next.ID, &next.Username, &next.Firstname, &next.Lastname, &next.Tasks, &next.DoneTasks, &next.Issues, &next.ClosedIssues); err != nil {
			return err
		}
		perUser = append(perUser, next)
		return nil
	}, request.Project, request.StartDate, request.DueDate)
	return
}

func (db *PostgresConnection) TimeEfficiencyReport(ctx context.Context, request reqresp.ReportRequest) (perUser []generic.TimeEfficiencyReport, err error) {
	perUser = []generic.TimeEfficiencyReport{}
	err = db.genericQueryForeach(ctx, "tasks_report", func(r *sql.Rows) error {
		var next generic.TimeEfficiencyReport
		if err := r.Scan(&next.ID, &next.Username, &next.Firstname, &next.Lastname, &next.TrackedRecords, &next.SummaryHours, &next.AveragePerDay, &next.AveragePerMonth); err != nil {
			return err
		}
		next.SummaryHours /= 3600.0
		next.AveragePerDay /= 3600.0
		next.AveragePerMonth /= 3600.0
		perUser = append(perUser, next)
		return nil
	}, request.Project, request.StartDate, request.DueDate)
	return
}
