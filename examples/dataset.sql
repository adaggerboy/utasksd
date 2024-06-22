
create user bob with encrypted password '11111111'; grant generic to bob;
create user charlie with encrypted password '11111111'; grant generic to charlie;
create user dana with encrypted password '11111111'; grant generic to dana;
create user eveline with encrypted password '11111111'; grant generic to eveline;
create user bismark with encrypted password '11111111'; grant generic to bismark;
create user matizont with encrypted password '11111111'; grant generic to matizont;
create user supporter with encrypted password '11111111'; grant generic to supporter;
create user ireness with encrypted password '11111111'; grant generic to ireness;
create user jacksone with encrypted password '11111111'; grant generic to jacksone;

SET session_replication_role = replica;

INSERT INTO projects (owner_id, name, logo, description)
VALUES (1, 'Project Phoenix', 'null', 'Revamping the internal system.');


-- Workers
INSERT INTO users (email, username, firstname, lastname, avatar, is_active)
VALUES 
('bob@example.com', 'bob', 'Bob', 'Smith', 'null', true),
('charlie@example.com', 'charlie', 'Charlie', 'Brown', 'null', true),
('dana@example.com', 'dana', 'Dana', 'White', 'null', true),
('eve@example.com', 'eveline', 'Eve', 'Black', 'null', true);

-- Managers
INSERT INTO users (email, username, firstname, lastname, avatar, is_active)
VALUES 
('fiona@example.com', 'bismark', 'Fiona', 'Green', 'null', true),
('george@example.com', 'matizont', 'George', 'Blue', 'null', true);

-- Support
INSERT INTO users (email, username, firstname, lastname, avatar, is_active)
VALUES 
('hank@example.com', 'supporter', 'Hank', 'Grey', 'null', true),
('irene@example.com', 'ireness', 'Irene', 'Red', 'null', true),
('jack@example.com', 'jacksone', 'Jack', 'Yellow', 'null', true);

-- Workers
INSERT INTO users_projects (user_id, project_id, access_level)
VALUES 
(2, 1, 'worker'),
(3, 1, 'worker'),
(4, 1, 'worker'),
(5, 1, 'worker');

-- Managers
INSERT INTO users_projects (user_id, project_id, access_level)
VALUES 
(6, 1, 'manager'),
(7, 1, 'manager'),
(1, 1, 'owner');

-- Support
INSERT INTO users_projects (user_id, project_id, access_level)
VALUES 
(8, 1, 'support'),
(9, 1, 'support'),
(10, 1, 'support');

INSERT INTO tasks (assigner_id, project_id, name, description, status, priority, start_date, due_date)
VALUES 
(6, 1, 'Set up infrastructure', 'Initial setup of project infrastructure.', 'in-progress', 'high', '2024-06-01', '2024-07-01'),
(7, 1, 'Design UI Mockups', 'Create mockups for the new UI.', 'done', 'middle', '2024-06-05', '2024-07-10'),
(6, 1, 'Database Schema', 'Design and implement the database schema.', 'done', 'high', '2024-06-10', '2024-07-15'),
(7, 1, 'Implement Auth System', 'Develop user authentication and authorization.', 'done', 'high', '2024-06-15', '2024-08-01'),
(6, 1, 'Setup CI/CD', 'Configure continuous integration and deployment.', 'done', 'middle', '2024-06-20', '2024-07-20'),
(7, 1, 'Backend API', 'Develop backend APIs for the application.', 'in-progress', 'high', '2024-06-25', '2024-08-10'),
(6, 1, 'Frontend Integration', 'Integrate frontend with backend APIs.', 'to-do', 'middle', '2024-07-01', '2024-08-15'),
(7, 1, 'Testing', 'Perform unit and integration testing.', 'to-do', 'middle', '2024-07-05', '2024-08-20'),
(6, 1, 'User Feedback', 'Gather and analyze user feedback.', 'done', 'low', '2024-07-10', '2024-08-25'),
(7, 1, 'Bug Fixing', 'Fix reported bugs and issues.', 'to-do', 'high', '2024-07-15', '2024-08-30'),
(6, 1, 'Deployment', 'Deploy the application to production.', 'to-do', 'high', '2024-07-20', '2024-09-01'),
(7, 1, 'Post-Deployment Support', 'Provide support post-deployment.', 'to-do', 'low', '2024-07-25', '2024-09-05'),
(6, 1, 'Documentation', 'Create comprehensive documentation.', 'to-do', 'middle', '2024-07-30', '2024-09-10'),
(7, 1, 'Performance Optimization', 'Optimize the application for performance.', 'to-do', 'high', '2024-08-01', '2024-09-15'),
(6, 1, 'Security Review', 'Conduct a thorough security review.', 'to-do', 'high', '2024-08-05', '2024-09-20'),
(7, 1, 'UX Enhancements', 'Enhance the user experience based on feedback.', 'to-do', 'middle', '2024-08-10', '2024-09-25'),
(6, 1, 'Accessibility Improvements', 'Ensure the application is accessible.', 'to-do', 'middle', '2024-08-15', '2024-09-30'),
(7, 1, 'SEO Optimization', 'Optimize the application for search engines.', 'to-do', 'low', '2024-08-20', '2024-10-05'),
(6, 1, 'Feature X Development', 'Develop the new feature X.', 'in-progress', 'high', '2024-08-25', '2024-10-10'),
(7, 1, 'Feature Y Development', 'Develop the new feature Y.', 'in-progress', 'high', '2024-08-30', '2024-10-15');


INSERT INTO issues (project_id, publisher_id, name, reporter, description, status)
VALUES 
(1, 8, 'Issue 1', 'Bob', 'Issue 1 description.', 'open'),
(1, 9, 'Issue 2', 'Charlie', 'Issue 2 description.', 'closed'),
(1, 8, 'Issue 3', 'Dana', 'Issue 3 description.', 'closed'),
(1, 8, 'Issue 4', 'Eve', 'Issue 4 description.', 'closed'),
(1, 8, 'Issue 5', 'Fiona', 'Issue 5 description.', 'closed'),
(1, 8, 'Issue 6', 'George', 'Issue 6 description.', 'open'),
(1, 8, 'Issue 7', 'Hank', 'Issue 7 description.', 'open'),
(1, 9, 'Issue 8', 'Irene', 'Issue 8 description.', 'open'),
(1, 10, 'Issue 9', 'Jack', 'Issue 9 description.', 'closed'),
(1, 9, 'Issue 10', 'Bob', 'Issue 10 description.', 'open');

INSERT INTO tasks_issues (task_id, issue_id)
VALUES 
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5),
(6, 6),
(7, 7),
(8, 8),
(9, 9),
(10, 10);

-- Assign each task to a worker
INSERT INTO tasks_assignees (task_id, assignee_id)
VALUES 
(1, 2),
(2, 3),
(3, 4),
(4, 5),
(5, 2),
(6, 3),
(7, 4),
(8, 5),
(9, 2),
(10, 3),
(11, 4),
(12, 5),
(13, 2),
(14, 3),
(15, 4),
(16, 5),
(17, 2),
(18, 3),
(19, 4),
(20, 5);

-- Generate tracking records for the tasks
INSERT INTO tracking_records (task_id, text, seconds, end_date)
VALUES 
-- Task 1
(1, 'Initial setup of infrastructure', 7200, '2024-06-02 10:00:00'),
(1, 'Configuring network settings', 3600, '2024-06-03 14:00:00'),
(1, 'Setting up servers', 5400, '2024-06-04 09:00:00'),

-- Task 2
(2, 'Creating UI mockups for homepage', 5400, '2024-06-06 15:00:00'),
(2, 'Reviewing UI mockups', 1800, '2024-06-07 11:00:00'),

-- Task 3
(3, 'Database schema design', 10800, '2024-06-11 12:00:00'),
(3, 'Normalizing tables', 7200, '2024-06-12 16:00:00'),

-- Task 4
(4, 'Developing authentication system', 9600, '2024-06-16 18:00:00'),
(4, 'Implementing JWT tokens', 5400, '2024-06-17 14:00:00'),
(4, 'Testing authentication flow', 3600, '2024-06-18 09:00:00'),

-- Task 5
(5, 'Setting up CI/CD pipeline', 6000, '2024-06-21 14:00:00'),
(5, 'Configuring Jenkins', 7200, '2024-06-22 11:00:00'),

-- Task 6
(6, 'Backend API development', 8500, '2024-06-26 16:00:00'),
(6, 'Implementing REST endpoints', 5400, '2024-06-27 14:00:00'),
(6, 'Writing API documentation', 3600, '2024-06-28 09:00:00'),

-- Task 7
(7, 'Frontend integration', 7000, '2024-07-02 11:00:00'),
(7, 'Connecting to backend APIs', 5400, '2024-07-03 15:00:00'),

-- Task 8
(8, 'Unit testing', 4000, '2024-07-06 09:00:00'),
(8, 'Writing test cases', 3600, '2024-07-07 10:00:00'),
(8, 'Running tests', 3200, '2024-07-08 13:00:00'),

-- Task 9
(9, 'Analyzing user feedback', 3600, '2024-07-11 10:00:00'),
(9, 'Compiling feedback report', 5400, '2024-07-12 14:00:00'),

-- Task 10
(10, 'Bug fixing session', 7200, '2024-07-16 17:00:00'),
(10, 'Debugging reported issues', 5400, '2024-07-17 11:00:00'),
(10, 'Testing fixes', 3600, '2024-07-18 16:00:00'),

-- Task 11
(11, 'Deployment preparation', 6500, '2024-07-21 13:00:00'),
(11, 'Finalizing deployment scripts', 3600, '2024-07-22 10:00:00'),

-- Task 12
(12, 'Providing post-deployment support', 3000, '2024-07-26 19:00:00'),
(12, 'Monitoring application', 3600, '2024-07-27 14:00:00'),

-- Task 13
(13, 'Writing documentation', 7200, '2024-07-31 08:00:00'),
(13, 'Reviewing documentation', 3600, '2024-08-01 09:00:00'),

-- Task 14
(14, 'Optimizing performance', 8600, '2024-08-02 14:00:00'),
(14, 'Profiling application', 5400, '2024-08-03 11:00:00'),

-- Task 15
(15, 'Conducting security review', 7200, '2024-08-05 14:00:00'),
(15, 'Implementing security fixes', 5400, '2024-08-06 09:00:00'),
(15, 'Testing security measures', 3600, '2024-08-07 16:00:00');