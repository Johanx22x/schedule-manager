package main

import (
    "fmt"
    "flag"
    "os"
    "io/ioutil"
    "gopkg.in/yaml.v2"

    "github.com/Johanx22x/schedule-manager/ansi"
)

const (
    // Constant DIR is the directory where the folder with the courses is located.
    DIR = "/home/johanw/university"

    // Constant CURRENT_COURSE is the name of the folder that contains the current course.
    CURRENT_COURSE = DIR + "/current-course"
)

var (
    // Variable currentCourse is the name of the folder that will be linked to the 
    // current course. 
    currentCourseFlag string

    // Variable newCourse is the name of the folder that will be created.
    newCourseFlag string

    // Variable listCourses is a boolean that is true if the flag -lc is used.
    listCoursesFlag bool

    // Variable help is a boolean that is true if the flag -h is used.
    helpFlag bool
)

func init() {
    // Flag -cc is used to change the simbolic link to the current course. It takes 
    // the name of folder as argument.
    flag.StringVar(&currentCourseFlag, "cc", "", "Change the current course")

    // Flag -nc is used to create a new course. It takes the name of the folder as
    // argument.
    flag.StringVar(&newCourseFlag, "nc", "", "Create a new course")

    // Flag -lc is used to list all the courses.
    flag.BoolVar(&listCoursesFlag, "lc", false, "List all the courses")

    // Flag -h is used to print the help message.
    flag.BoolVar(&helpFlag, "h", false, "Print the help message")

    // Parse the flags.
    flag.Parse()
}

// Function getFolders returns the names of the folders into the directory dir.
func getFolders(dir string) []string {
    // Use the os package to get the files into the directory dir.
    files, err := os.ReadDir(dir)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    // Create a slice of strings to store the names of the folders.
    folders := make([]string, 0)

    // Iterate over the files into the directory dir.
    for _, file := range files {
        // If the file is a directory, add its name to the slice folders.
        if file.IsDir() {
            folders = append(folders, file.Name())
        }
    }

    // Return the slice folders.
    return folders
}

// Function currentCourse returns the name of the folder that is linked to the 
// current course.
func currentCourse() string {
    // Use ioutil.ReadFIle to read the content of `info.yaml`.
    data, err := ioutil.ReadFile(CURRENT_COURSE + "/info.yaml")
    if err != nil {
        fmt.Println(err)
        return ""
    }

    // Create a map to store the content of `info.yaml`.
    info := make(map[string]string)

    // Use yaml.Unmarshal to unmarshal the content of `info.yaml` into the map info.
    err = yaml.Unmarshal(data, &info)
    if err != nil {
        fmt.Println(err)
        return ""
    }

    // Change the spaces in the name of the course with dashes.
    course := info["title"]
    for i := 0; i < len(course); i++ {
        if course[i] == ' ' {
            course = course[:i] + "-" + course[i+1:]
        }
    }

    // Return the name of the folder that is linked to the current course.
    return course
}

// Function listCourses lists all the courses.
func listCourses() {
    currentCourse := currentCourse()
    // Iterate over the folders into the subfolders of the directory DIR.
    for _, semester := range getFolders(DIR) {
        for _, course := range getFolders(DIR + "/" + semester) {
            // If the folder is the current course, print it in green.
            if course == currentCourse {
                fmt.Println(semester + " -> " + ansi.Green + course + ansi.Reset + " (current course)")
            } else {
                fmt.Println(semester + " -> " + course)
            }
        }
    }
}

// Function changeCurrentCourse changes the current course.
func changeCurrentCourse(name string) {
    // Map of valid courses, store the course name as key and the semester as value.
    validCourses := make(map[string]string)

    // Obtain the valid courses.
    for _, semester := range getFolders(DIR) {
        for _, course := range getFolders(DIR + "/" + semester) {
            validCourses[course] = semester
        }
    }

    // Check if the course is valid.
    if _, ok := validCourses[name]; !ok {
        fmt.Println("The course " + name + " is not valid")
        fmt.Println("Use the flag -lc to list all the courses")
        return 
    }
    
    // Change the current course.
    err := os.Remove(CURRENT_COURSE)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = os.Symlink(DIR + "/" + validCourses[name] + "/" + name, CURRENT_COURSE)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("The current course has been changed to " + name)
}

// Function help prints the help message.
func help() {
    fmt.Println("Usage: schedule-manager [options]")
    fmt.Println("Options:")
    fmt.Println("  -cc <course>  Change the current course")
    fmt.Println("  -nc <course>  Create a new course")
    fmt.Println("  -lc           List all the courses")
    fmt.Println("  -h            Print the help message")
}

func main() {
    // If the flag -h is used, print the help message.
    if helpFlag {
        help()
        return
    }

    // If the flag -lc is used, list all the courses.
    if listCoursesFlag {
        fmt.Println("Courses:")
        listCourses()
        return
    }

    // If the flag -cc is used, change the current course.
    if currentCourseFlag != "" {
        changeCurrentCourse(currentCourseFlag)
        return
    }

    // If the flag -nc is used, create a new course.
    if newCourseFlag != "" {
        fmt.Println("Create a new course")
        return
    }

    // If no flag is used, print the help message.
    help()
}
