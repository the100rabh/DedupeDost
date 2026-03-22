#!/bin/bash

# ============================================
# DupDel Cleanup Script
# ============================================
# Generated: 2026-03-18 14:05:33
# Source Directory: /dup-del
# Total files to delete: 28
# Total space to recover: 162.09 KB
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

# Group 1: employees_old.csv (262 B)
delete_file "/dup-del/test_data/documents/employees.csv"
delete_file "/dup-del/test_data/work/shared/employees.csv"
# Keeping: /dup-del/test_data/archive/2023/employees_old.csv

# Group 2: deploy_backup.sh (197 B)
delete_file "/dup-del/test_data/work/project_a/deploy.sh"
delete_file "/dup-del/test_data/work/project_b/deploy.sh"
# Keeping: /dup-del/test_data/backup/monthly/deploy_backup.sh

# Group 3: data_backup.json (302 B)
delete_file "/dup-del/test_data/documents/data.json"
delete_file "/dup-del/test_data/work/shared/data.json"
# Keeping: /dup-del/test_data/backup/weekly/data_backup.json

# Group 4: app.log (329 B)
delete_file "/dup-del/test_data/backup/daily/app.log"
delete_file "/dup-del/test_data/temp/cache/app.log"
delete_file "/dup-del/test_data/temp/logs/app.log"
# Keeping: /dup-del/test_data/archive/2023/app.log

# Group 5: README.md (183 B)
delete_file "/dup-del/test_data/documents/projects/README.md"
delete_file "/dup-del/test_data/work/project_b/README.md"
# Keeping: /dup-del/test_data/work/project_a/README.md

# Group 6: file1.bin (10.00 KB)
delete_file "/dup-del/test_data/backup/daily/file1.bin"
delete_file "/dup-del/test_data/temp/cache/file1.bin"
# Keeping: /dup-del/test_data/downloads/file1.bin

# Group 7: .keep (0 B)
delete_file "/dup-del/test_data/temp/empty2.txt"
delete_file "/dup-del/test_data/temp/empty1.txt"
# Keeping: /dup-del/test_data/temp/cache/.keep

# Group 8: archive.zip (50.00 KB)
delete_file "/dup-del/test_data/downloads/archive.zip"
# Keeping: /dup-del/test_data/backup/weekly/archive.zip

# Group 9: meeting_old.txt (237 B)
delete_file "/dup-del/test_data/backup/daily/meeting_backup.txt"
delete_file "/dup-del/test_data/documents/notes/meeting.txt"
delete_file "/dup-del/test_data/personal/work/meeting.txt"
delete_file "/dup-del/test_data/work/shared/meeting_notes.txt"
# Keeping: /dup-del/test_data/archive/old/meeting_old.txt

# Group 10: report_backup.txt (134 B)
delete_file "/dup-del/test_data/backup/daily/report.txt"
delete_file "/dup-del/test_data/documents/report.txt"
# Keeping: /dup-del/test_data/archive/2023/report_backup.txt

# Group 11: index.html (170 B)
delete_file "/dup-del/test_data/work/project_a/index.html"
# Keeping: /dup-del/test_data/documents/index.html

# Group 12: large_file.txt (43.84 KB)
delete_file "/dup-del/test_data/backup/daily/large_file.txt"
delete_file "/dup-del/test_data/documents/large_file.txt"
# Keeping: /dup-del/test_data/archive/old/large_file.txt

# Group 13: config_backup.ini (85 B)
delete_file "/dup-del/test_data/work/project_a/config.ini"
delete_file "/dup-del/test_data/documents/config.ini"
delete_file "/dup-del/test_data/work/project_b/config.ini"
# Keeping: /dup-del/test_data/backup/weekly/config_backup.ini

# ============================================
# Cleanup Summary
# ============================================
# Duplicate groups processed: 13
# Files marked for deletion: 28
# Space to be recovered: 162.09 KB
#
# To execute this script:
#   1. Review the files to be deleted
#   2. Run: DRY_RUN=true ./cleanup.sh (preview)
#   3. Run: ./cleanup.sh (actual deletion)
# ============================================
