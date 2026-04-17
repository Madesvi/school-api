CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY, 
    first_name VARCHAR(255) NOT NULL, 
    last_name VARCHAR(255) NOT NULL, 
    email VARCHAR(255) NOT NULL UNIQUE,
    teacher_id INTEGER NOT NULL, 
    CONSTRAINT fk_student_teacher 
        FOREIGN KEY (teacher_id) 
        REFERENCES teachers(id) 
        ON DELETE CASCADE
);

ALTER SEQUENCE IF EXISTS students_id_seq RESTART WITH 100;
