#!/bin/bash

# =============================================================================
# DedupeDost - Test File Generator
# =============================================================================
# This script creates a test directory structure with duplicate and unique files
# for testing the DedupeDost duplicate file scanner application.
#
# Usage: ./gen_test_files.sh [test_directory]
# Example: ./gen_test_files.sh ./test_data
# =============================================================================

TEST_DIR="${1:-./test_data}"

echo "╔══════════════════════════════════════════════════════════╗"
echo "║       DedupeDost Test File Generator                         ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Clean up existing test directory
if [ -d "$TEST_DIR" ]; then
    echo "Removing existing test directory: $TEST_DIR"
    rm -rf "$TEST_DIR"
fi

echo "Creating test directory structure at: $TEST_DIR"

# Create main directory structure
mkdir -p "$TEST_DIR"/{documents,photos,videos,downloads,backup,archive,work,personal,temp}
mkdir -p "$TEST_DIR/documents"/{reports,notes,projects}
mkdir -p "$TEST_DIR/photos"/{2023,2024,vacation}
mkdir -p "$TEST_DIR/videos"/{tutorials,personal}
mkdir -p "$TEST_DIR/downloads"/{installers,archives}
mkdir -p "$TEST_DIR/backup"/{daily,weekly,monthly}
mkdir -p "$TEST_DIR/archive"/{2022,2023,old}
mkdir -p "$TEST_DIR/work"/{project_a,project_b,shared}
mkdir -p "$TEST_DIR/personal"/{finance,health,hobbies}
mkdir -p "$TEST_DIR/temp"/{cache,logs}

echo "Creating duplicate files..."

# Duplicate Set 1: Text documents (3 copies)
cat > "$TEST_DIR/documents/report.txt" << 'EOF'
This is an important document.
It contains critical information.
Please do not delete this file.
Created: 2024-01-15
Author: John Doe
EOF
cp "$TEST_DIR/documents/report.txt" "$TEST_DIR/backup/daily/report.txt"
cp "$TEST_DIR/documents/report.txt" "$TEST_DIR/archive/2023/report_backup.txt"

# Duplicate Set 2: Configuration files (4 copies)
cat > "$TEST_DIR/documents/config.ini" << 'EOF'
[settings]
theme = dark
language = en
auto_save = true
max_files = 1000
timeout = 30
EOF
cp "$TEST_DIR/documents/config.ini" "$TEST_DIR/work/project_a/config.ini"
cp "$TEST_DIR/documents/config.ini" "$TEST_DIR/work/project_b/config.ini"
cp "$TEST_DIR/documents/config.ini" "$TEST_DIR/backup/weekly/config_backup.ini"

# Duplicate Set 3: Notes (5 copies)
cat > "$TEST_DIR/documents/notes/meeting.txt" << 'EOF'
Meeting Notes - 2024-03-15
Attendees: Alice, Bob, Charlie
Topics:
1. Project timeline review
2. Budget allocation
3. Resource planning
Action items:
- Alice: Prepare presentation
- Bob: Update documentation
- Charlie: Schedule follow-up
EOF
cp "$TEST_DIR/documents/notes/meeting.txt" "$TEST_DIR/work/shared/meeting_notes.txt"
mkdir -p "$TEST_DIR/personal/work"
cp "$TEST_DIR/documents/notes/meeting.txt" "$TEST_DIR/personal/work/meeting.txt"
cp "$TEST_DIR/documents/notes/meeting.txt" "$TEST_DIR/backup/daily/meeting_backup.txt"
cp "$TEST_DIR/documents/notes/meeting.txt" "$TEST_DIR/archive/old/meeting_old.txt"

# Duplicate Set 4: README files (3 copies)
cat > "$TEST_DIR/work/project_a/README.md" << 'EOF'
# Project README

## Description
This is a sample project for testing.

## Installation
```
npm install
npm start
```

## Usage
Run the application and enjoy!

## License
MIT License
EOF
cp "$TEST_DIR/work/project_a/README.md" "$TEST_DIR/work/project_b/README.md"
cp "$TEST_DIR/work/project_a/README.md" "$TEST_DIR/documents/projects/README.md"

# Duplicate Set 5: Log files (4 copies)
cat > "$TEST_DIR/temp/logs/app.log" << 'EOF'
[2024-03-15 10:30:00] INFO: Application started
[2024-03-15 10:30:01] INFO: Loading configuration
[2024-03-15 10:30:02] INFO: Database connected
[2024-03-15 10:30:03] WARNING: Cache miss
[2024-03-15 10:30:04] INFO: Processing request
[2024-03-15 10:30:05] INFO: Request completed
[2024-03-15 10:30:06] INFO: Application shutdown
EOF
cp "$TEST_DIR/temp/logs/app.log" "$TEST_DIR/temp/cache/app.log"
cp "$TEST_DIR/temp/logs/app.log" "$TEST_DIR/backup/daily/app.log"
cp "$TEST_DIR/temp/logs/app.log" "$TEST_DIR/archive/2023/app.log"

# Duplicate Set 6: JSON data files (3 copies)
cat > "$TEST_DIR/documents/data.json" << 'EOF'
{
    "users": [
        {"id": 1, "name": "Alice", "email": "alice@example.com"},
        {"id": 2, "name": "Bob", "email": "bob@example.com"},
        {"id": 3, "name": "Charlie", "email": "charlie@example.com"}
    ],
    "settings": {
        "theme": "dark",
        "notifications": true
    }
}
EOF
cp "$TEST_DIR/documents/data.json" "$TEST_DIR/work/shared/data.json"
cp "$TEST_DIR/documents/data.json" "$TEST_DIR/backup/weekly/data_backup.json"

# Duplicate Set 7: Shell scripts (3 copies)
cat > "$TEST_DIR/work/project_a/deploy.sh" << 'EOF'
#!/bin/bash
# Sample deployment script
echo "Starting deployment..."
echo "Checking dependencies..."
echo "Installing packages..."
echo "Configuring services..."
echo "Deployment complete!"
exit 0
EOF
cp "$TEST_DIR/work/project_a/deploy.sh" "$TEST_DIR/work/project_b/deploy.sh"
cp "$TEST_DIR/work/project_a/deploy.sh" "$TEST_DIR/backup/monthly/deploy_backup.sh"

# Duplicate Set 8: CSV files (3 copies)
cat > "$TEST_DIR/documents/employees.csv" << 'EOF'
id,name,email,department,salary
1,John Doe,john@example.com,Engineering,75000
2,Jane Smith,jane@example.com,Marketing,65000
3,Bob Johnson,bob@example.com,Sales,70000
4,Alice Brown,alice@example.com,HR,60000
5,Charlie Wilson,charlie@example.com,Engineering,80000
EOF
cp "$TEST_DIR/documents/employees.csv" "$TEST_DIR/work/shared/employees.csv"
cp "$TEST_DIR/documents/employees.csv" "$TEST_DIR/archive/2023/employees_old.csv"

# Duplicate Set 9: HTML files (2 copies)
cat > "$TEST_DIR/documents/index.html" << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>Welcome</h1>
    <p>This is a test page.</p>
</body>
</html>
EOF
cp "$TEST_DIR/documents/index.html" "$TEST_DIR/work/project_a/index.html"

# Duplicate Set 10: Large text files (3 copies)
for i in {1..1000}; do echo "Line $i - Test data for duplicate detection"; done > "$TEST_DIR/documents/large_file.txt"
cp "$TEST_DIR/documents/large_file.txt" "$TEST_DIR/backup/daily/large_file.txt"
cp "$TEST_DIR/documents/large_file.txt" "$TEST_DIR/archive/old/large_file.txt"

# Duplicate Set 11: Binary files (3 copies)
echo "Creating binary duplicate files..."
dd if=/dev/urandom of="$TEST_DIR/downloads/file1.bin" bs=1024 count=10 2>/dev/null
cp "$TEST_DIR/downloads/file1.bin" "$TEST_DIR/backup/daily/file1.bin"
cp "$TEST_DIR/downloads/file1.bin" "$TEST_DIR/temp/cache/file1.bin"

# Duplicate Set 12: Medium binary files (2 copies)
dd if=/dev/urandom of="$TEST_DIR/downloads/archive.zip" bs=1024 count=50 2>/dev/null
cp "$TEST_DIR/downloads/archive.zip" "$TEST_DIR/backup/weekly/archive.zip"

# Duplicate Set 13: PNG image files (3 copies)
echo "Creating duplicate image files..."
# Create a simple PNG file (1x1 red pixel using base64)
echo 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg==' | base64 -d > "$TEST_DIR/photos/image.png"
cp "$TEST_DIR/photos/image.png" "$TEST_DIR/photos/2024/image_copy.png"
cp "$TEST_DIR/photos/image.png" "$TEST_DIR/backup/daily/image_backup.png"

# Duplicate Set 14: JPEG image files (2 copies)
# Create a minimal valid JPEG (1x1 pixel)
echo '/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAn/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBEQCEAwEPwAAf/9k=' | base64 -d > "$TEST_DIR/photos/photo.jpg"
cp "$TEST_DIR/photos/photo.jpg" "$TEST_DIR/photos/vacation/photo_vacation.jpg"

# Duplicate Set 15: GIF image files (2 copies)
# Create a simple GIF file (1x1 pixel)
echo 'R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7' | base64 -d > "$TEST_DIR/photos/animation.gif"
cp "$TEST_DIR/photos/animation.gif" "$TEST_DIR/photos/2023/animation_old.gif"

# Duplicate Set 16: BMP image files (2 copies)
# Create a simple BMP file (2x2 pixels, 24-bit)
printf 'BM6\x00\x00\x00\x00\x00\x00\x006\x00\x00\x00(\x00\x00\x00\x02\x00\x00\x00\x02\x00\x00\x00\x01\x00\x18\x00\x00\x00\x00\x000\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xff\xff\xff\x00\xff\xff\xff\x00\xff\xff\xff\x00\xff\xff\xff\x00' > "$TEST_DIR/photos/bitmap.bmp"
cp "$TEST_DIR/photos/bitmap.bmp" "$TEST_DIR/backup/weekly/bitmap_backup.bmp"

# Duplicate Set 17: Video files (3 copies) - MP4
echo "Creating duplicate video files..."
# Generate a small valid MP4 video (2 seconds, solid color)
ffmpeg -y -f lavfi -i color=c=blue:s=320x240:d=2 -c:v libx264 -pix_fmt yuv420p "$TEST_DIR/videos/vacation.mp4" >/dev/null 2>&1
cp "$TEST_DIR/videos/vacation.mp4" "$TEST_DIR/videos/vacation_copy.mp4"
cp "$TEST_DIR/videos/vacation.mp4" "$TEST_DIR/backup/daily/vacation_backup.mp4"

# Duplicate Set 18: Video files (2 copies) - AVI
# Generate a small valid AVI video (2 seconds, solid red color)
ffmpeg -y -f lavfi -i color=c=red:s=320x240:d=2 -c:v mjpeg "$TEST_DIR/videos/presentation.avi" >/dev/null 2>&1
cp "$TEST_DIR/videos/presentation.avi" "$TEST_DIR/videos/presentation_old.avi"

echo "Creating unique files..."

# Unique documents
echo "This is unique document 1 with specific content." > "$TEST_DIR/documents/unique1.txt"
echo "This is unique document 2 with different content." > "$TEST_DIR/documents/unique2.txt"
echo "This is unique document 3 with yet another content." > "$TEST_DIR/documents/unique3.txt"

# Unique notes
echo -e "Shopping List:\n- Milk\n- Eggs\n- Bread\n- Butter" > "$TEST_DIR/documents/notes/shopping_list.txt"
echo -e "Password hints (not actual passwords!):\n- Email: pet name + birth year\n- Bank: street number + zip" > "$TEST_DIR/documents/notes/passwords.txt"
echo -e "Project Ideas:\n1. Mobile app for task management\n2. Web scraper for price comparison\n3. Chatbot for customer support" > "$TEST_DIR/documents/notes/ideas.txt"

# Unique work files
echo -e "Project A Requirements:\n- Python 3.8+\n- Flask\n- PostgreSQL\n- Redis" > "$TEST_DIR/work/project_a/requirements.txt"
echo -e "Project B Specifications:\n- Node.js 16+\n- Express\n- MongoDB\n- Docker" > "$TEST_DIR/work/project_b/specs.txt"
echo -e "Team Contacts:\n- Alice: alice@company.com\n- Bob: bob@company.com\n- Charlie: charlie@company.com" > "$TEST_DIR/work/shared/team_contacts.txt"

# Unique personal files
echo -e "2024 Budget:\nIncome: \$5000/month\nRent: \$1500\nFood: \$500\nTransport: \$300\nSavings: \$1500" > "$TEST_DIR/personal/finance/budget_2024.txt"
echo -e "Workout Plan:\nMonday: Chest & Triceps\nWednesday: Back & Biceps\nFriday: Legs & Shoulders" > "$TEST_DIR/personal/health/workout_plan.txt"
echo -e "Reading List:\n- Clean Code by Robert Martin\n- The Pragmatic Programmer\n- Design Patterns" > "$TEST_DIR/personal/hobbies/reading_list.txt"

# Unique temp files
echo "Temporary cache data - can be deleted" > "$TEST_DIR/temp/cache/temp_data.tmp"
echo -e "Error log with unique entries\n[2024-03-15] Error: Connection timeout\n[2024-03-16] Error: File not found" > "$TEST_DIR/temp/logs/error.log"

# Unique binary files
echo "Creating unique binary files..."
dd if=/dev/urandom of="$TEST_DIR/downloads/unique1.dat" bs=1024 count=2 2>/dev/null
dd if=/dev/urandom of="$TEST_DIR/downloads/unique2.dat" bs=1024 count=4 2>/dev/null
dd if=/dev/urandom of="$TEST_DIR/downloads/unique3.dat" bs=1024 count=8 2>/dev/null
dd if=/dev/urandom of="$TEST_DIR/temp/cache/cache1.dat" bs=1024 count=1 2>/dev/null
dd if=/dev/urandom of="$TEST_DIR/temp/cache/cache2.dat" bs=1024 count=2 2>/dev/null

echo "Creating hidden files..."
echo -e "*.log\n*.tmp\n*.cache\nnode_modules/" > "$TEST_DIR/.gitignore"
echo "These are secret notes." > "$TEST_DIR/documents/.hidden_notes.txt"
echo -e "BACKUP_DIR=/mnt/backup\nRETENTION_DAYS=30" > "$TEST_DIR/backup/.backup_config"
echo -e "DB_HOST=localhost\nDB_USER=admin\nDB_PASS=secret123" > "$TEST_DIR/.env"

echo "Creating files with special names..."
echo "This is a file with 'copy' in name" > "$TEST_DIR/documents/file (copy).txt"
echo "This is a file with multiple 'copy' in name" > "$TEST_DIR/documents/file (copy) (copy).txt"
echo "This is another copy variation" > "$TEST_DIR/backup/file_copy.txt"
echo "This is a backup variation" > "$TEST_DIR/archive/file_backup.txt"

echo "Creating empty files..."
touch "$TEST_DIR/temp/empty1.txt"
touch "$TEST_DIR/temp/empty2.txt"
touch "$TEST_DIR/temp/cache/.keep"

echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║              Test Files Generation Complete              ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""
echo "Test Directory: $TEST_DIR"
echo ""

# Count files
TOTAL_FILES=$(find "$TEST_DIR" -type f | wc -l)
TOTAL_DIRS=$(find "$TEST_DIR" -type d | wc -l)

echo "Statistics:"
echo "  ────────────────────────────────────────"
echo "  Total Files Created:    $TOTAL_FILES"
echo "  Total Directories:      $TOTAL_DIRS"
echo "  ────────────────────────────────────────"
echo ""
echo "Duplicate Sets Created:"
echo "  ────────────────────────────────────────"
echo "  1.  report.txt          (3 copies)"
echo "  2.  config.ini          (4 copies)"
echo "  3.  meeting.txt         (5 copies)"
echo "  4.  README.md           (3 copies)"
echo "  5.  app.log             (4 copies)"
echo "  6.  data.json           (3 copies)"
echo "  7.  deploy.sh           (3 copies)"
echo "  8.  employees.csv       (3 copies)"
echo "  9.  index.html          (2 copies)"
echo "  10. large_file.txt      (3 copies)"
echo "  11. file1.bin           (3 copies)"
echo "  12. archive.zip         (2 copies)"
echo "  13. image.png           (3 copies) [PNG]"
echo "  14. photo.jpg           (2 copies) [JPEG]"
echo "  15. animation.gif       (2 copies) [GIF]"
echo "  16. bitmap.bmp          (2 copies) [BMP]"
echo "  17. vacation.mp4        (3 copies) [MP4 Video]"
echo "  18. presentation.avi    (2 copies) [AVI Video]"
echo ""
echo "Next Steps:"
echo "  ────────────────────────────────────────"
echo "  1. Run DedupeDost CLI:  ./dedupedost -dir $TEST_DIR"
echo "  2. Run DedupeDost GUI:  ./dedupedost-gui"
echo "  3. Review detected duplicates"
echo "  4. Generate and review cleanup script"
echo ""
echo "✓ Test files ready for DedupeDost testing!"
