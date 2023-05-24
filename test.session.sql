-- @block
CREATE TABLE IF NOT EXISTS Quiz(
    id INT PRIMARY KEY AUTO_INCREMENT,
    quiz_title VARCHAR(255) NOT NULL,
    owner_id INTEGER NOT NULL
);

-- @block
INSERT INTO Quiz (quiz_title, owner_id) VALUES ('Quiz 1', 1);

-- @block
SELECT * FROM Quiz;

-- @block
CREATE TABLE IF NOT EXISTS Sections(
    id INT AUTO_INCREMENT,
    section_title VARCHAR(255) NOT NULL,
    section_background VARCHAR(255) NOT NULL,
    quiz_id INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (quiz_id) REFERENCES Quiz(id)
);

-- @block
INSERT INTO Sections (section_title, section_background, quiz_id) VALUES 
('Section 1', 'https://www.w3schools.com/css/img_fjords.jpg', 1);

-- @block
SELECT * FROM quiz
INNER JOIN sections on  quiz_id = quiz.id;
