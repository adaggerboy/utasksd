do $$
begin
    if not exists (select 1 from pg_type where typname = 'user_access_level') then
        create type user_access_level as enum ('worker', 'support', 'manager', 'owner');
    end if;
    if not exists (select 1 from pg_type where typname = 'task_priority') then
        create type task_priority as enum ('lowest', 'low', 'middle', 'high', 'highest');
    end if;
    if not exists (select 1 from pg_type where typname = 'task_status') then
        create type task_status as enum ('to-do', 'in-progress', 'testing', 'waiting', 'done');
    end if;
    if not exists (select 1 from pg_type where typname = 'issue_status') then
        create type issue_status as enum ('open', 'closed', 'reopened');
    end if;
    if not exists (select 1 from pg_type where typname = 'task_dependency_type') then
        create type task_dependency_type as enum ('blocked-by', 'includes');
    end if;
    if not exists (select 1 from pg_roles where rolname = 'administrator') then
        create role administrator;
    end if;
    if not exists (select 1 from pg_roles where rolname = 'director') then
        create role director;
    end if;
    if not exists (select 1 from pg_roles where rolname = 'generic') then
        create role generic;
    end if;
end $$;



create table if not exists
    attachments (
        cdn_path varchar(255) not null,
        primary key (cdn_path)
    );

insert into attachments (cdn_path) values ('null') on conflict do nothing;

create table if not exists
    users (
        id serial,
        email varchar(40) not null unique check (email ~* '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
        username varchar(40) not null unique check (username ~* '^[a-zA-Z]+$'),
        firstname varchar(30) not null check (firstname ~* '^[a-zA-Z]+$'),
        lastname varchar(30) not null check (lastname ~* '^[a-zA-Z]+$'),
        avatar varchar(255) not null,
        is_active bool not null default true,
        is_admin bool not null default false,
        is_director bool not null default false,
        created_at timestamp not null default now (),
        primary key (id),
        foreign key (avatar) references attachments (cdn_path) deferrable
    );


create table if not exists
    projects (
        id serial,
        owner_id integer not null,
        name varchar(60) not null,
        logo varchar(255) not null,
        description varchar(500),
        created_at timestamp not null default now (),
        primary key (id),
        foreign key (owner_id) references users (id) deferrable,
        foreign key (logo) references attachments (cdn_path) deferrable
    );


create table if not exists
    users_projects (
        user_id integer not null,
        project_id integer not null,
        access_level user_access_level not null default 'worker',
        primary key (user_id, project_id),
        foreign key (user_id) references users (id) deferrable,
        foreign key (project_id) references projects (id) on delete cascade deferrable
    );

create table if not exists
    issues (
        id serial,
        project_id integer not null,
        publisher_id integer not null,
        name varchar(80) not null,
        reporter varchar(80),
        description varchar(5000) not null,
        status issue_status not null default 'open',
        open_date timestamp not null default now (),
        close_date timestamp default null,
        primary key (id),
        foreign key (project_id) references projects (id) on delete cascade deferrable,
        foreign key (publisher_id) references users (id) deferrable
    );

create table if not exists
    issues_attachments (
        issue_id integer not null,
        attachment varchar(255) not null,
        primary key (issue_id, attachment),
        foreign key (issue_id) references issues (id) on delete cascade deferrable,
        foreign key (attachment) references attachments (cdn_path) deferrable
    );


create table if not exists
    tasks (
        id serial,
        assigner_id integer not null,
        project_id integer not null,
        name varchar(80) not null,
        description varchar(8192) not null,
        status task_status not null default 'to-do',
        priority task_priority not null default 'middle',
        start_date timestamp not null default now (),
        due_date timestamp,
        primary key (id),
        foreign key (assigner_id) references users (id) deferrable,
        foreign key (project_id) references projects (id) on delete cascade deferrable
    );

create table if not exists
    tasks_attachments (
        task_id integer not null,
        attachment varchar(255) not null,
        primary key (task_id, attachment),
        foreign key (task_id) references tasks (id) on delete cascade deferrable,
        foreign key (attachment) references attachments (cdn_path) deferrable
    );


create table if not exists
    tasks_assignees (
        task_id integer not null unique,
        assignee_id integer not null,
        primary key (task_id, assignee_id),
        foreign key (task_id) references tasks (id) on delete cascade deferrable,
        foreign key (assignee_id) references users (id) deferrable
    );


create table if not exists
    tasks_issues (
        task_id integer not null,
        issue_id integer not null,
        primary key (task_id, issue_id),
        foreign key (task_id) references tasks (id) on delete cascade deferrable,
        foreign key (issue_id) references issues (id) on delete cascade deferrable
    );


create table if not exists
    tasks_dependencies (
        task_id integer not null,
        dep_task_id integer not null,
        dependency_type task_dependency_type not null,
        primary key (task_id, dep_task_id),
        foreign key (task_id) references tasks (id) on delete cascade deferrable,
        foreign key (dep_task_id) references tasks (id) on delete cascade deferrable
    );

create table if not exists
    comments (
        id serial,
        task_id integer not null,
        user_id integer not null,
        text varchar(500) not null,
        date timestamp not null default now (),
        primary key (id),
        foreign key (task_id) references tasks (id) on delete cascade deferrable,
        foreign key (user_id) references users (id)
    );

create table if not exists
    tracking_records (
        id serial,
        task_id integer not null,
        text varchar(100) not null,
        seconds integer not null check (seconds > 0),
        end_date timestamp not null default now (),
        primary key (id),
        foreign key (task_id) references tasks (id) on delete cascade deferrable
    );



create or replace function handle_project_changes() returns trigger as $$
begin
    delete from users_projects where project_id = new.id and access_level = 'owner';
    insert into users_projects (user_id, project_id, access_level)
    values (new.owner_id, new.id, 'owner');
    return new;
end;
$$ language plpgsql;

drop trigger if exists project_changes_trigger on projects;
create trigger project_changes_trigger 
after insert or update on projects 
for each row execute function handle_project_changes();

create or replace function prevent_dependent_task_creation()
returns trigger as $$
begin
    if exists (select 1 from tasks t1 join tasks t2 on t1.project_id != t2.project_id where t1.id = new.task_id and t2.id = new.dep_task_id) then
        raise exception 'dependent task must belong to the same project as the task';
    end if;
    if exists (select 1 from tasks_dependencies where task_id = new.dep_task_id) then
        raise exception 'dependent task already has dependent tasks';
    end if;
    if exists (select 1 from tasks_dependencies where dep_task_id = new.task_id) then
        raise exception 'dependent task cannot be reflexive';
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists prevent_dependent_task_creation_trigger on tasks_dependencies;
create trigger prevent_dependent_task_creation_trigger
before insert on tasks_dependencies
for each row execute function prevent_dependent_task_creation();

create or replace function check_dependent_tasks_exist()
returns trigger as $$
begin
    if exists (
        select 1
        from tasks_dependencies
        where task_id = old.id
    ) then
        raise exception 'cannot delete task because it has dependent tasks';
    end if;
    return old;
end;
$$ language plpgsql;

drop trigger if exists refuse_task_deletion on tasks;
create trigger refuse_task_deletion
before delete on tasks
for each row
execute function check_dependent_tasks_exist();

create or replace function check_dependent_tasks_status()
returns trigger as $$
begin
    if exists (
        select 1
        from tasks_dependencies td
        join tasks t on td.dep_task_id = t.id
        where td.task_id = new.id
        and t.status != 'done'
    ) then
        raise exception 'cannot change task status to "done" because there are dependent tasks with status not "done"';
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists refuse_done_status_change on tasks;
create trigger refuse_done_status_change
before update on tasks
for each row when (new.status = 'done') execute function check_dependent_tasks_status();

create or replace function update_close_date()
returns trigger as $$
begin
    if new.status = 'closed' then
        new.close_date := now();
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists update_close_date_trigger on issues;
create trigger update_close_date_trigger
before update on issues
for each row
when (old.status is distinct from new.status)
execute function update_close_date();

create or replace function prevent_linking_to_other_projects()
returns trigger as $$
begin
    if not exists (
        select 1
        from tasks, issues
        where tasks.id = new.task_id
        and issues.id = new.issue_id
        and tasks.project_id = issues.project_id
    ) then
    raise exception 'task cannot be linked to an issue from another project';
    end if;
    return new;
end;
$$ language plpgsql;

drop trigger if exists check_link_to_other_projects on tasks_issues;
create trigger check_link_to_other_projects
before insert or update on tasks_issues
for each row
execute function prevent_linking_to_other_projects();

create or replace function prevent_assigning_to_other_projects()
returns trigger as $$
begin
    if exists (
        select 1
        from tasks
        inner join users_projects up on up.project_id = tasks.project_id
        where tasks.id = new.task_id
        and up.user_id = new.assignee_id
    ) then
        return new;
    else
        raise exception 'task cannot be assigned to a user from another project';
    end if;
end;
$$ language plpgsql;

drop trigger if exists check_assign_to_other_projects on tasks_assignees;
create trigger check_assign_to_other_projects
before insert or update on tasks_assignees
for each row
execute function prevent_assigning_to_other_projects();


grant all privileges on table attachments to administrator;
grant insert, select on table attachments to generic;

grant select, insert, update on table users to administrator;
grant select, update (username, firstname, lastname, email, avatar, is_active) on table users to director;
grant select, update (username, firstname, lastname, email, avatar, is_active) on table users to public;

grant all privileges on table projects to administrator;
grant all privileges on table projects to director;
grant select on table projects to generic;

grant all privileges on table users_projects to administrator;
grant all privileges on table users_projects to director;
grant select on table users_projects to generic;

grant all privileges on table issues to generic;
grant select, insert, delete on table issues_attachments to generic;
grant all privileges on table tasks to generic;
grant select, insert, delete on table tasks_attachments to generic;
grant select, insert, delete on table tasks_assignees to generic;
grant select, insert, delete on table tasks_issues to generic;
grant all privileges on table tasks_dependencies to generic;
grant select, insert, delete on table comments to generic;
grant select, insert, delete on table tracking_records to generic;



create or replace function tasks_report(_project_id int, _start_date timestamp, _due_date timestamp)
returns table (
    user_id numeric,
    user_name text,
    user_firstname text,
    user_lastname text,
    num_tasks numeric,
    done_tasks numeric,
    num_issues numeric,
    closed_issues numeric
) as
$$
begin
    return query 
    with task_user_stats as (
        select
            tasks_assignees.assignee_id,
            tasks_assignees.task_id,
            tasks.status,
            count (issue_id) issues,
            count (case when issues.status = 'closed' then 1 end) done_issues
        from tasks_assignees
        inner join tasks on tasks_assignees.task_id = tasks.id
        left join tasks_issues on tasks_assignees.task_id = tasks_issues.task_id
        left join issues on issues.id = tasks_issues.issue_id
        where 
            tasks.project_id = coalesce(_project_id, tasks.project_id) and
            start_date >= coalesce(_start_date, start_date) and
            due_date <= coalesce(_due_date, due_date)
        group by tasks_assignees.task_id, tasks_assignees.assignee_id, tasks.id
    ) select 
        users.id::numeric as user_id,
        users.username::text as user_name,
        users.firstname::text as user_firstname,
        users.lastname::text as user_lastname,
        count(task_id)::numeric as num_tasks,
        count(case when status = 'done' then 1 end)::numeric as done_tasks,
        sum(issues)::numeric as num_issues,
        sum(done_issues)::numeric as closed_issues
    from task_user_stats
    inner join users on task_user_stats.assignee_id = users.id
    group by users.id;
end;
$$
language plpgsql;
grant execute on function tasks_report to director;


create or replace function time_efficiency_report(_project_id int, _start_date timestamp, _due_date timestamp)
returns table (
    user_id numeric,
    user_name text,
    user_firstname text,
    user_lastname text,
    tracked_steps numeric,
    summary numeric,
    average_per_day numeric,
    average_per_month numeric
) as
$$
begin
    return query 
    with track_user_stats as (
        select
            tracking_records.seconds, 
            sum(tracking_records.seconds) over (partition by tasks.project_id, tasks_assignees.assignee_id, date_trunc('day', tracking_records.end_date)) as sum_day, 
            sum(tracking_records.seconds) over (partition by tasks.project_id, tasks_assignees.assignee_id, date_trunc('month', tracking_records.end_date)) as sum_month, 
            tasks_assignees.assignee_id 
        from tracking_records 
        inner join tasks_assignees on tasks_assignees.task_id = tracking_records.task_id 
        inner join tasks on tasks.id = tasks_assignees.task_id
        where 
            tasks.project_id = coalesce(_project_id, tasks.project_id) and
            tracking_records.end_date >= coalesce(_start_date, tracking_records.end_date) and
            tracking_records.end_date <= coalesce(_due_date, tracking_records.end_date)
    ) 
    select 
        users.id::numeric as user_id,
        users.username::text as user_name,
        users.firstname::text as user_firstname,
        users.lastname::text as user_lastname,
        count(track_user_stats.seconds)::numeric as tracked_steps,
        sum(track_user_stats.seconds)::numeric as summary,
        avg(track_user_stats.sum_day)::numeric as average_per_day,
        avg(track_user_stats.sum_month)::numeric as average_per_month
    from track_user_stats
    inner join users on track_user_stats.assignee_id = users.id
    group by users.id;
end;
$$
language plpgsql;
grant execute on function time_efficiency_report to director;

create or replace procedure delete_unused_attachments()
language plpgsql
as $$
declare
    attachment_record record;
begin
    for attachment_record in 
        select cdn_path
        from attachments a
        where not exists (select 1 from users u where u.avatar = a.cdn_path)
        and not exists (select 1 from projects p where p.logo = a.cdn_path)
        and not exists (select 1 from issues_attachments ia where ia.attachment = a.cdn_path)
        and not exists (select 1 from tasks_attachments ta where ta.attachment = a.cdn_path)
        and not a.cdn_path = 'null' 
    loop
        delete from attachments where cdn_path = attachment_record.cdn_path;
    end loop;
end;
$$;
grant execute on procedure delete_unused_attachments to administrator;



grant select, usage on all sequences in schema public to generic;
