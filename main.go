package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
    "time"
	"gopkg.in/yaml.v2"

	"github.com/Johanx22x/schedule-manager/ansi"
)

const (
    // PROGRAM_DIR is the directory where the program logs the errors.
    PROGRAM_DIR = "/home/johanw/.schedule-manager"

    // Log file.
    LOG = PROGRAM_DIR + "/log.txt"

    // Constant DIR is the directory where the folder with the courses is located.
    DIR = "/home/johanw/university"

    // Constant CURRENT_COURSE is the name of the folder that contains the current course.
    CURRENT_COURSE = DIR + "/current-course"
)

var (
    // Variable currentCourse is the name of the folder that will be linked to the 
    // current course. 
    currentCourseFlag string

    // Variable listCourses is a boolean that is true if the flag -lc is used.
    listCoursesFlag bool

    // Variable showPDF is a boolean that is true if the flag -sPdf 
    showPDFFlag bool

    // Variable openc is a boolean that is true if the flag -oc is used.
    opencFlag bool

    // Variable openLink is a boolean that is true if the flag -cl is used.
    openLinkFlag bool

    // Variable getCourseName is a boolean that is true if the flag -cn is used.
    getCourseNameFlag bool

    // Variable daemon is a boolean that is true if the flag -d is used.
    daemonFlag bool

    // Variable showCourses is a boolean that is true if the flag -sc is used.
    showCoursesFlag bool

    // Variable help is a boolean that is true if the flag -h is used.
    helpFlag bool
)

func init() {
    // Create the directory PROGRAM_DIR.
    err := os.MkdirAll(PROGRAM_DIR, 0755)
    if err != nil {
        log.Println(err)
        return
    }

    // Create the log file. 
    file, err := os.OpenFile(LOG, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Println(err)
        return 
    }

    // Use the log package to log the errors.
    log.SetOutput(file)

    // Flag -cc is used to change the simbolic link to the current course. It takes 
    // the name of folder as argument.
    flag.StringVar(&currentCourseFlag, "cc", "", "Change the current course")

    // Flag -lc is used to list all the courses.
    flag.BoolVar(&listCoursesFlag, "lc", false, "List all the courses")

    // Flag -sPdf is used to show the pdf of the current course.
    flag.BoolVar(&showPDFFlag, "sPdf", false, "Show the pdf of the current course")

    // Flag -oc is used to open the folder of the current course.
    flag.BoolVar(&opencFlag, "oc", false, "Open the folder of the current course")

    // Flag -cl is used to open the link of the current course.
    flag.BoolVar(&openLinkFlag, "cl", false, "Open the link of the current course")

    // Flag -cn is used to get the name of the current course.
    flag.BoolVar(&getCourseNameFlag, "cn", false, "Get the name of the current course")

    // Flag -d is used to run the program as a daemon.
    flag.BoolVar(&daemonFlag, "d", false, "Run the program as a daemon (not implemented yet)")

    // Flag -sc is used to show the courses.
    flag.BoolVar(&showCoursesFlag, "sc", false, "Show the courses")

    // Flag -h is used to print the help message.
    flag.BoolVar(&helpFlag, "h", false, "Print the help message")

    // Parse the flags.
    flag.Parse()
}

// Function that prints the available courses.
func showCourses() {
    // Iterate over the folders into the directory DIR.
    for _, semester := range getFolders(DIR) {
        for _, course := range getFolders(DIR + "/" + semester) {
            fmt.Println(course)
        }
    }
}

// Function getFolders returns the names of the folders into the directory dir.
func getFolders(dir string) []string {
    // Use the os package to get the files into the directory dir.
    files, err := os.ReadDir(dir)
    if err != nil {
        log.Println(err)
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
        log.Println(err)
        return ""
    }

    // Create a map to store the content of `info.yaml`.
    info := make(map[string]string)

    // Use yaml.Unmarshal to unmarshal the content of `info.yaml` into the map info.
    err = yaml.Unmarshal(data, &info)
    if err != nil {
        log.Println(err)
        return ""
    }

    // Change the spaces in the name of the course with dashes.
    course := info["title"]

    // Return the name of the folder that is linked to the current course.
    return course
}

// Function listCourses lists all the courses.
func listCourses() {
    currentCourse := currentCourse()
    for i := 0; i < len(currentCourse); i++ {
        if currentCourse[i] == ' ' {
            currentCourse = currentCourse[:i] + "-" + currentCourse[i+1:]
        }
    }

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
        fmt.Println(ansi.Red + "The course " + name + " is an invalid course" + ansi.Reset)
        fmt.Println("Use the flag -lc to list all the courses")

        log.Println("The course " + name + " is an invalid course")
        return 
    }
    
    // Change the current course.
    err := os.Remove(CURRENT_COURSE)
    if err != nil {
        log.Println(err)
        return
    }

    err = os.Symlink(DIR + "/" + validCourses[name] + "/" + name, CURRENT_COURSE)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println("The current course has been changed to " + name)
}

// Function showPDF shows the pdf of the current course.
func showPDF() {
    // Search the pdf file into CURRENT_COURSE.
    files, err := os.ReadDir(CURRENT_COURSE)
    if err != nil {
        log.Println(err)
        return 
    }

    // Store the name of the pdf file.
    var pdf string

    // Iterate over the files into CURRENT_COURSE.
    for _, file := range files {
        // If the file is a pdf file, store its name into pdf.
        if filepath.Ext(file.Name()) == ".pdf" {
            pdf = file.Name()
        }
    }

    // Search it into /build.
    if pdf == "" {
        files, err = os.ReadDir(CURRENT_COURSE + "/build")
        if err != nil {
            log.Println(err)
            return 
        }

        // Iterate over the files into CURRENT_COURSE/build.
        for _, file := range files {
            // If the file is a pdf file, store its name into pdf.
            if filepath.Ext(file.Name()) == ".pdf" {
                pdf = file.Name()
            }
        }
    }

    // If pdf is empty, the pdf file has not been found.
    if pdf == "" {
        fmt.Println(ansi.Red + "The pdf file has not been found" + ansi.Reset)
        log.Println("The pdf file has not been found")
        return 
    }

    // Use the exec package to open the pdf file.
    err = exec.Command("zathura", CURRENT_COURSE + "/" + pdf).Run()
    if err != nil {
        log.Println(err)
        return 
    }
}

// Function openLink opens the link of the current course.
func openLink() {
    // Use ioutil.ReadFIle to read the content of `info.yaml`.
    data, err := ioutil.ReadFile(CURRENT_COURSE + "/info.yaml")
    if err != nil {
        log.Println(err)
        return
    }

    // Create a map to store the content of `info.yaml`.
    info := make(map[string]string)

    // Use yaml.Unmarshal to unmarshal the content of `info.yaml` into the map info.
    err = yaml.Unmarshal(data, &info)
    if err != nil {
        log.Println(err)
        return
    }

    // Use the exec package to open the link.
    err = exec.Command("firefox", info["link"]).Run()
    if err != nil {
        log.Println(err)
        return
    }
}

// Function help prints the help message.
func help() {
    fmt.Println("Usage: schedule-manager [options]")
    fmt.Println("Options:")
    fmt.Println("  -cc <course>  Change the current course")
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

    // If the flag -sPdf is used, show the pdf of the current course.
    if showPDFFlag {
        showPDF()
        return 
    }

    // If the flag -oc is used, open the current course.
    if opencFlag {
        // Use the exec package to open the current course.
        err := exec.Command("alacritty", "-e", "nvim", CURRENT_COURSE).Run()
        if err != nil {
            log.Println(err)
            return 
        }
        return 
    }

    // If the flag -cl is used, open the current course link in the browser.
    if openLinkFlag {
        openLink()
        return 
    }

    // If the flag -cn is used, print the name of the current course.
    if getCourseNameFlag {
        if daemonFlag {
            // If the flag -d is used, print the name of the current course every 5 seconds.
            for {
                fmt.Println(currentCourse())
                time.Sleep(5 * time.Second)
            }
        } else {
            // If the flag -d is not used, print the name of the current course.
            fmt.Println(currentCourse())
        }
        return 
    }

    // If the flag -sc is used, show the courses.
    if showCoursesFlag {
        showCourses()
        return 
    }

    // If no flag is used, print the help message.
    log.Println("No flag used")
    help()
}
