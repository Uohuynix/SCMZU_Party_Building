package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// RestoreFromFile restores database from SQL file
func RestoreFromFile(sqlDB *sql.DB, sqlFile string) error {
	// Check if file exists
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		return fmt.Errorf("SQL file does not exist: %s", sqlFile)
	}

	// Read SQL file
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Split SQL statements by semicolon
	sqlContent := string(sqlBytes)
	statements := strings.Split(sqlContent, ";")

	// Execute each SQL statement
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err = sqlDB.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to execute SQL: %v", err)
		}
	}

	return nil
}

// CreateTables creates database tables
func CreateTables(sqlDB *sql.DB) error {
	createTablesSQL := `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin','teacher','student') DEFAULT 'student',
    name VARCHAR(50),
    branch VARCHAR(100),
    group_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_no VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    gender ENUM('male','female'),
    ethnicity VARCHAR(50),
    birth_date DATE NULL,
    education VARCHAR(50),
    phone VARCHAR(20),
    id_card VARCHAR(20),
    major_class VARCHAR(100),
    type ENUM('masses','activist','candidate','probationary','formal') DEFAULT 'masses',
    branch VARCHAR(100),
    group_name VARCHAR(100),
    photo_url VARCHAR(255),
    admission_date DATE NULL,
    conversion_date DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS trainings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    period VARCHAR(50),
    student_id INT,
    unit VARCHAR(100),
    start_date DATE NULL,
    end_date DATE NULL,
    score ENUM('excellent','good','pass','fail'),
    certificate_no VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS developments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT,
    join_league_date VARCHAR(50) NULL,
    apply_date VARCHAR(50) NULL,
    party_talk_date VARCHAR(50) NULL,
    league_recommend_date VARCHAR(50) NULL,
    activist_date VARCHAR(50) NULL,
    candidate_date VARCHAR(50) NULL,
    superior_report_date VARCHAR(50) NULL,
    committee_review_date VARCHAR(50) NULL,
    party_committee_pre_review_date VARCHAR(50) NULL,
    branch_meeting_date VARCHAR(50) NULL,
    conversion_apply_date VARCHAR(50) NULL,
    conversion_branch_meeting_date VARCHAR(50) NULL,
    probation_date VARCHAR(50) NULL,
    conversion_date VARCHAR(50) NULL,
    transfer_date VARCHAR(50) NULL,
    introduction_date VARCHAR(50) NULL,
    status ENUM('masses','activist','candidate','probationary','formal') DEFAULT 'masses',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS rewards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT,
    type ENUM('reward','punishment') NOT NULL,
    content VARCHAR(255) NOT NULL,
    date DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id)
);
`

	_, err := sqlDB.Exec(createTablesSQL)
	return err
}

// InsertSampleData inserts sample data
func InsertSampleData(sqlDB *sql.DB) error {
	sampleDataSQL := `
INSERT IGNORE INTO users (username, password, role, name, branch, group_name) VALUES 
-- 系统管理员：不隶属于任何具体党支部/党小组
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'Administrator', NULL, NULL),
-- 教师：作为“计算机学院2021级本科生党支部”的管理员，可管理本支部下所有小组
('teacher1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'teacher', 'Teacher Zhang', '计算机学院2021级本科生党支部', NULL),
-- 学生：隶属于具体党支部和党小组，便于教师/管理员按组织管理
('student1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', 'Student Li', '计算机学院2021级本科生党支部', '第一小组');

INSERT IGNORE INTO students (student_no, name, gender, ethnicity, birth_date, education, phone, id_card, major_class, type, branch, group_name, admission_date) VALUES 
('2021001', 'Zhang San', 'male', 'Han', '2000-01-01', 'Bachelor', '13800138001', '110101200001011234', 'Computer Science', 'masses', '计算机学院2021级本科生党支部', '第一小组', '2021-09-01'),
('2021002', 'Li Si', 'female', 'Han', '2000-02-02', 'Bachelor', '13800138002', '110101200002021234', 'Software Engineering', 'activist', '计算机学院2021级本科生党支部', '第一小组', '2021-09-01'),
('2021003', 'Wang Wu', 'male', 'Han', '2000-03-03', 'Bachelor', '13800138003', '110101200003031234', 'Network Engineering', 'candidate', '计算机学院2021级本科生党支部', '第二小组', '2021-09-01');

INSERT IGNORE INTO trainings (period, student_id, unit, start_date, end_date, score, certificate_no) VALUES 
('2021 Spring', 1, 'Marxism School', '2021-03-01', '2021-06-30', 'excellent', 'CERT2021001'),
('2021 Spring', 2, 'Marxism School', '2021-03-01', '2021-06-30', 'good', 'CERT2021002'),
('2021 Fall', 3, 'Marxism School', '2021-09-01', '2021-12-31', 'pass', 'CERT2021003');

INSERT IGNORE INTO developments (student_id, apply_date, activist_date, status) VALUES 
(1, '2021-10-01', '2021-11-01', 'masses'),
(2, '2021-09-01', '2021-10-01', 'activist'),
(3, '2021-08-01', '2021-09-01', 'candidate');

INSERT IGNORE INTO rewards (student_id, type, content, date) VALUES 
(1, 'reward', 'Excellent Student Leader', '2021-12-01'),
(2, 'reward', 'Three Good Student', '2021-12-01'),
(3, 'reward', 'Active Learner', '2021-12-01');
`

	_, err := sqlDB.Exec(sampleDataSQL)
	return err
}
