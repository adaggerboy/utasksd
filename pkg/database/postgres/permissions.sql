


alter table attachments enable row level security;

alter table users enable row level security;

alter table user_credentials enable row level security;

alter table projects enable row level security;

alter table users_projects enable row level security;

alter table issues enable row level security;

alter table issues_attachments enable row level security;

alter table tasks enable row level security;

alter table tasks_attachments enable row level security;

alter table tasks_dependencies enable row level security;

alter table tasks_issues enable row level security;

alter table tasks_dependencies enable row level security;

alter table comments enable row level security;

alter table tracking_records enable row level security;


create role director;
create role administrator;

grant administrator to authuser;

grant select on users to public;

grant all on schema public to administrator;
grant all on all sequences in schema public TO administrator;

create policy users_self_update_policy on users
for update
using (username = current_user);

-- project

grant insert on projects to director;

create policy projects_update_policy on projects
for update, delete
using (owner_id = (select id from users where username = current_user));

create policy projects_select_policy on projects
for select
using (
    exists (
        select 1 from users_projects
        where users_projects.project_id = projects.id
        and users_projects.user_id = (select id from users where username = current_user)
    )
);



create policy users_projects_insert_update_policy on users_projects
for all
to director
using (
    exists (
        select 1 from projects
        where projects.id = users_projects.project_id
        and projects.owner_id = (select id from users where username = current_user)
    )
);

create policy users_projects_select_policy on users_projects
for select
using (
    exists (
        select 1 from users_projects 
        where users_projects.user_id = (select id from users where username = current_user)
    )
);




create policy issues_insert_policy on issues
for insert
using (
    exists (
        select 1 from users_projects up
        where up.project_id = issues.project_id
        and up.user_id = (select id from users where username = current_user)
        and (up.access_level = 'support' or up.access_level = 'owner')
    )
);

create policy issues_update_policy on issues
for update
using (
    (publisher_id = (select id from users where username = current_user))
    or
    exists (
        select 1 from projects
        where projects.id = issues.project_id
        and projects.owner_id = (select id from users where username = current_user)
    )
)
with check (
    and project_id = old.project_id
    and publisher_id = old.publisher_id
    and reporter = old.reporter
    and open_date = old.open_date
);

create policy issues_select_policy on issues
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = issues.project_id
        and up.user_id = (select id from users where username = current_user)
    )
);






create policy tasks_insert_policy on tasks
for insert
using (
    exists (
        select 1 from projects
        where projects.id = tasks.project_id
        and (
            projects.owner_id = (select id from users where username = current_user)
            or exists (
                select 1 from users_projects
                where users_projects.project_id = tasks.project_id
                and users_projects.user_id = (select id from users where username = current_user)
                and users_projects.access_level = 'manager'
            )
        )
    )
);

create policy tasks_update_policy on tasks
for update
using (
    exists (
        select 1 from projects
        where projects.id = tasks.project_id
        and projects.owner_id = (select id from users where username = current_user)
    )
)
with check (
    (project_id = old.project_id)
);

create policy tasks_update_assigner_policy on tasks
for update
using (
    (assigner_id = (select id from users where username = current_user))
)
with check (
    (project_id = old.project_id) and (assigner_id = old.assigner_id)
);

create policy tasks_select_policy on tasks
for select
using (
    exists (
        select 1 from users_projects
        where users_projects.project_id = tasks.project_id
        and users_projects.user_id = (select id from users where username = current_user)
    )
);











create policy tasks_attachments_insert_update_delete_policy on tasks_attachments
for insert, update, delete
using (
    exists (
        select 1 from tasks
        where tasks.id = tasks_attachments.task_id
        and (
            tasks.assigner_id = (select id from users where username = current_user)
            or exists (
                select 1 from projects
                where projects.id = tasks.project_id
                and projects.owner_id = (select id from users where username = current_user)
            )
        )
    )
);

create policy tasks_attachments_select_policy on tasks_attachments
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tasks_attachments.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);



-- combination tables

create policy issues_attachments_insert_update_delete_policy on issues_attachments
for insert, update, delete
using (
    exists (
        select 1 from issues
        where issues.id = issues_attachments.issue_id
        and issues.publisher_id = (select id from users where username = current_user)
    )
);

create policy issues_attachments_select_policy on issues_attachments
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from issues
            where issues.id = issues_attachments.issue_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);

create policy tasks_attachments_insert_update_delete_policy on tasks_attachments
for insert, update, delete
using (
    exists (
        select 1 from tasks
        where tasks.id = tasks_attachments.task_id
        and tasks.assigner_id = (select id from users where username = current_user)
    )
);

create policy tasks_attachments_select_policy on tasks_attachments
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tasks_attachments.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);

create policy tasks_issues_insert_update_delete_policy on tasks_issues
for insert, update, delete
using (
    exists (
        select 1 from tasks
        where tasks.id = tasks_issues.task_id
        and tasks.assigner_id = (select id from users where username = current_user)
    )
);

create policy tasks_issues_select_policy on tasks_issues
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tasks_issues.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);


create policy tasks_assignees_insert_update_delete_policy on tasks_assignees
for insert, update, delete
using (
    exists (
        select 1 from tasks
        where tasks.id = tasks_assignees.task_id
        and tasks.assigner_id = (select id from users where username = current_user)
    )
);

create policy tasks_assignees_select_policy on tasks_assignees
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tasks_assignees.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);


create policy tasks_dependencies_insert_update_delete_policy on tasks_dependencies
for insert, update, delete
using (
    exists (
        select 1 from tasks
        where tasks.id = tasks_dependencies.task_id
        and tasks.assigner_id = (select id from users where username = current_user)
    )
);

create policy tasks_dependencies_select_policy on tasks_dependencies
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tasks_dependencies.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);


create policy tracking_records_insert_update_delete_policy on tracking_records
for insert, delete
using (
    exists (
        select 1 from tasks_assignees
        where tasks_assignees.task_id = tracking_records.task_id
        and tasks_assignees.assignee_id = (select id from users where username = current_user)
    )
);

create policy tracking_records_select_policy on tracking_records
for select
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = tracking_records.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);


create policy comments_insert_update_delete_policy on comments
for update, delete
using (
    user_id = (select id from users where username = current_user)
) with check (
    (user_id = old.user_id)
);

create policy comments_select_policy on comments
for select, insert
using (
    exists (
        select 1 from users_projects up
        where up.project_id = (
            select project_id from tasks
            where tasks.id = comments.task_id
        )
        and up.user_id = (select id from users where username = current_user)
    )
);









