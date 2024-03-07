package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	//file_with_path := "go run main.go /home/dinesh//"
	//o build main.go -o /home/dinesh/untarz

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		fmt.Println("Provide single argument")
		fmt.Println("./untarz rootfolder")
		os.Exit(1)
	}


	//untarzst(file_with_path)
	//root := "/home/dinesh/"
	root := argsWithoutProg[0]

	fmt.Println("Make sure sudo apt-get install zstd")
	fmt.Println("Root folder is " + root)
	src_recipes := filepath.Join(root, "recipes")
	// bin_packages := filepath.Join(root, "packages")
	src_dst := filepath.Join(root, "recipes_srcs")
	// bin_dst := filepath.Join(root, "bin_target_pkgs")

	createDir(src_dst)
	// createDir(bin_dst)

	getSources(src_recipes, src_dst)
	// getSources(bin_packages, bin_dst)

}

func createDir(dst string) {
	//err := os.MkdirAll(dst, 777)

	isexist,_ := exists(dst)

	if isexist == true {
		err := os.RemoveAll(dst)
		if err != nil {
			fmt.Println("Error while removing "+dst)
			fmt.Println(err.Error())
			os.Exit(1)
		}

	}else {
		err := os.Mkdir(dst, 0755)
		if err != nil {
			fmt.Println("Error while Creating "+dst)
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}
func getSources(dir, destination_path string) {

	files, err := WalkDir(dir, []string{"tar.zst"})
	handleerr(err)
	fmt.Println(strconv.Itoa(len(files)) + " tar.zst files exists ===========!")
	//printFiles(files)
	for _, tarzstfile := range files {
		fmt.Println("Processing " + tarzstfile)
		untarzst(tarzstfile, destination_path)
	}
}
func printFiles(files []string) {
	i := 0
	for _, file := range files {
		i++
		fmt.Println(strconv.Itoa(i) + " " + file)
	}

}
func untarzst(tarzstfile string, destination_path string) {
	//destination_path := filepath.Join(filepath.Dir(tarzstfile), "src")

	// tar --use-compress-program=unzstd -xvf recipe-zlib.tar.zst --one-top-level -C destination_path
	cmd := exec.Command("tar", "--use-compress-program=unzstd", "-xvf", tarzstfile, "--one-top-level", "-C", destination_path)
	stdout, err := cmd.Output()
	handleerr(err)
	fmt.Println(string(stdout))
}
func handleerr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func WalkDir(root string, exts []string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		for _, s := range exts {
			if strings.HasSuffix(path, "."+s) {
				files = append(files, path)
				return nil
			}
		}
		return nil
	})
	return files, err
}
