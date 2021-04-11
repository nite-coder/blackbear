package config

// func TestMain(m *testing.M) {
// 	fmt.Println("run")
// 	// copy config file to executed file's directory
// 	path, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}
// 	srcPath := filepath.Join(path, "../../test/config/app_test.yml")

// 	path, err = os.Executable()
// 	if err != nil {
// 		panic(err)
// 	}
// 	dstPath := filepath.Join(filepath.Dir(path), "app.yml")
// 	iofile.CopyFile(srcPath, dstPath)

// 	fmt.Printf("dstpath: %s", dstPath)
// 	m.Run()
// 	fmt.Println("end")
// }
