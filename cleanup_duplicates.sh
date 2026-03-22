#!/bin/bash

# ============================================
# DedupeDost Cleanup Script
# ============================================
# Generated: 2026-03-22 16:15:11
# Source Directory: /home/the100rabh/code/personal/dedupedosteter/test_data/
# Total files to delete: 35
# Total space to recover: 208.28 KB
# ============================================


# DRY_RUN mode: set to "true" to preview without deleting
DRY_RUN=${DRY_RUN:-false}

# delete_file function - safely deletes a file
delete_file() {
    local file="$1"
    if [ "$DRY_RUN" = "true" ]; then
        echo "[DRY RUN] Would delete: $file"
    else
        if [ -f "$file" ]; then
            echo "Deleting: $file"
            rm -f "$file"
        else
            echo "Warning: File not found: $file"
        fi
    fi
}

# Group 1: presentation.avi (39.71 KB)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/videos/presentation_old.avi"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/videos/presentation.avi

# Group 2: bitmap.bmp (66 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/weekly/bitmap_backup.bmp"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/photos/bitmap.bmp

# Group 3: index.html (170 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/index.html"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/work/project_a/index.html

# Group 4: empty2.txt (0 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/temp/empty1.txt"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/temp/empty2.txt

# Group 5: employees.csv (262 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/archive/2023/employees_old.csv"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/shared/employees.csv"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/documents/employees.csv

# Group 6: file1.bin (10.00 KB)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/downloads/file1.bin"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/temp/cache/file1.bin"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/file1.bin

# Group 7: deploy.sh (197 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/monthly/deploy_backup.sh"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/project_b/deploy.sh"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/work/project_a/deploy.sh

# Group 8: image_backup.png (70 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/photos/image.png"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/photos/2024/image_copy.png"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/image_backup.png

# Group 9: animation_old.gif (42 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/photos/animation.gif"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/photos/2023/animation_old.gif

# Group 10: report_backup.txt (134 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/report.txt"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/report.txt"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/archive/2023/report_backup.txt

# Group 11: photo.jpg (284 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/photos/vacation/photo_vacation.jpg"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/photos/photo.jpg

# Group 12: large_file.txt (43.84 KB)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/large_file.txt"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/large_file.txt"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/archive/old/large_file.txt

# Group 13: app.log (329 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/temp/cache/app.log"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/app.log"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/temp/logs/app.log"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/archive/2023/app.log

# Group 14: config_backup.ini (85 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/config.ini"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/project_a/config.ini"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/project_b/config.ini"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/backup/weekly/config_backup.ini

# Group 15: meeting_old.txt (237 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/meeting_backup.txt"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/notes/meeting.txt"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/shared/meeting_notes.txt"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/personal/work/meeting.txt"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/archive/old/meeting_old.txt

# Group 16: README.md (183 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/documents/projects/README.md"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/project_b/README.md"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/work/project_a/README.md

# Group 17: vacation_backup.mp4 (2.98 KB)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/videos/vacation.mp4"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/videos/vacation_copy.mp4"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/backup/daily/vacation_backup.mp4

# Group 18: data.json (302 B)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/backup/weekly/data_backup.json"
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/work/shared/data.json"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/documents/data.json

# Group 19: archive.zip (50.00 KB)
delete_file "/home/the100rabh/code/personal/dedupedosteter/test_data/downloads/archive.zip"
# Keeping: /home/the100rabh/code/personal/dedupedosteter/test_data/backup/weekly/archive.zip

# ============================================
# Cleanup Summary
# ============================================
# Duplicate groups processed: 19
# Files marked for deletion: 35
# Space to be recovered: 208.28 KB
#
# To execute this script:
#   1. Review the files to be deleted
#   2. Run: DRY_RUN=true ./cleanup.sh (preview)
#   3. Run: ./cleanup.sh (actual deletion)
# ============================================
