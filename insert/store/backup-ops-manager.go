package store

//func (opts *nonRootRestartPodsFunc) checkOwnership(podName string) error {
//	var buf *bytes.Buffer
//	shSession := getCommand(podName, opts.db.Namespace, 3, buf)
//	output, err := shSession.Output()
//	if err != nil {
//		fmt.Printf("cant get output, err %s\n", err)
//		return err
//	}
//	errOutput := buf.String()
//	if errOutput != "" {
//		fmt.Printf("failed to execute command, stderr: %s", errOutput)
//		return err
//	}
//	if !isAllMongo(string(output)) {
//		return fmt.Errorf("user is not `mongodb` for all. %s \n", string(output))
//	}
//
//	shSession = getCommand(podName, opts.db.Namespace, 4, buf)
//	output, err = shSession.Output()
//	if err != nil {
//		fmt.Printf("cant get output, err %s\n", err)
//		return err
//	}
//	errOutput = buf.String()
//	if errOutput != "" {
//		fmt.Printf("failed to execute command, stderr: %s", errOutput)
//		return err
//	}
//	if !isAllMongo(string(output)) {
//		return fmt.Errorf("group is not `mongodb` for all. %s \n", string(output))
//	}
//	return nil
//}
//
//func getCommand(pod, ns string, column int, buf *bytes.Buffer) *sh.Session {
//	shell := sh.NewSession()
//	shell.ShowCMD = false
//	buf = &bytes.Buffer{}
//	shell.Stderr = buf
//
//	podName := fmt.Sprintf("pod/%s", pod)
//	kubectlCommand := []interface{}{
//		"exec", "-n", ns, podName, "-c", "mongodb", "--",
//	}
//	mgCommand := []interface{}{
//		"ls", "-l", "/data/db",
//	}
//	finalCommand := append(kubectlCommand, mgCommand...)
//
//	awkCommand := fmt.Sprintf("{print $%d}", column)
//	return shell.Command("kubectl", finalCommand...).Command("awk", awkCommand)
//}
//
//func isAllMongo(str string) bool {
//	klog.Infof("ssssssssssssssssssssssssssss %v\n", str)
//	slice := strings.Split(str, "\n")
//	for i := range slice {
//		if i == 0 {
//			continue
//		}
//		if slice[i] != "mongodb" {
//			return false
//		}
//	}
//	return true
//}
