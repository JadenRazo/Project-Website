@echo off
setlocal enabledelayedexpansion

REM Define the base directory
set "base=frontend\src"

REM Create directory arrays for our structure
set "dirs[0]=%base%\components\layout"
set "dirs[1]=%base%\components\projects"
set "dirs[2]=%base%\components\ui"
set "dirs[3]=%base%\components\forms"
set "dirs[4]=%base%\pages"
set "dirs[5]=%base%\hooks"
set "dirs[6]=%base%\context"
set "dirs[7]=%base%\api"
set "dirs[8]=%base%\types"
set "dirs[9]=%base%\utils"
set "dirs[10]=%base%\styles"
set "dirs[11]=%base%\assets\images"

REM Loop through and create each directory if it doesn't exist
for /L %%i in (0,1,11) do (
    if not exist "!dirs[%%i]!" (
        mkdir "!dirs[%%i]!"
        echo Created directory: !dirs[%%i]!
    ) else (
        echo Directory already exists: !dirs[%%i]!
    )
)

REM Create basic files if they don't exist
set "files[0]=%base%\components\layout\Navbar.tsx"
set "files[1]=%base%\components\layout\Footer.tsx"
set "files[2]=%base%\components\layout\Layout.tsx"
set "files[3]=%base%\components\projects\ProjectCard.tsx"
set "files[4]=%base%\components\projects\ProjectGrid.tsx"
set "files[5]=%base%\components\projects\ProjectDetails.tsx"
set "files[6]=%base%\pages\Home.tsx"
set "files[7]=%base%\pages\Projects.tsx"
set "files[8]=%base%\pages\About.tsx"
set "files[9]=%base%\pages\Contact.tsx"
set "files[10]=%base%\pages\ProjectDetail.tsx"
set "files[11]=%base%\pages\NotFound.tsx"
set "files[12]=%base%\api\client.ts"
set "files[13]=%base%\api\projects.ts"
set "files[14]=%base%\types\project.ts"
set "files[15]=%base%\hooks\useProjects.ts"
set "files[16]=%base%\App.tsx"
set "files[17]=%base%\index.tsx"
set "files[18]=%base%\router.tsx"
set "files[19]=%base%\styles\globals.css"

REM Loop through and create each file if it doesn't exist
for /L %%i in (0,1,19) do (
    if not exist "!files[%%i]!" (
        echo.> "!files[%%i]!"
        echo Created file: !files[%%i]!
    ) else (
        echo File already exists: !files[%%i]!
    )
)

echo.
echo Directory structure created successfully!