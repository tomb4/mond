package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"mond/wind/utils"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.PersistentFlags().StringVarP(&ServiceName, "name", "n", "", "service appid")
	serviceCmd.MarkPersistentFlagRequired("name")
	serviceCmd.PersistentFlags().Int32VarP(&Port, "port", "p", 0, "service port")
	serviceCmd.MarkPersistentFlagRequired("port")

	rootCmd.AddCommand(serviceCmd)
}

var (
	ServiceName string
	Port        int32
)

const (
	configYamlStr = `
appId: %s
port: %d
	`
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "mondc service -n=biz.demo -p=19001",
	Long:  "mondc service --name=biz.demo --port=19001",
	Run: func(cmd *cobra.Command, args []string) {
		// check dir
		serviceArr := strings.Split(ServiceName, ".")
		if len(serviceArr) != 2 {
			fmt.Println("Name must be split by . (e.g. biz.demo) ")
			return
		}
		serviceArr[0] = strings.ToLower(serviceArr[0])
		serviceArr[1] = strings.ToLower(serviceArr[1])
		AppId := utils.FirstUpper(serviceArr[0] + serviceArr[1])
		fmt.Println("AppId is", AppId)

		path, _ := os.Getwd()
		pathArr := strings.Split(path, "/")
		if pathArr[len(pathArr)-1] != "service" {
			fmt.Println("Please cd mond/service")
			return
		}

		folderPath := fmt.Sprintf("%s.%s", serviceArr[0], serviceArr[1])
		if utils.Exists(fmt.Sprintf("./%s", folderPath)) {
			fmt.Printf("dir %s alread existed \n", folderPath)
			return
		}
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}

		var (
			paramMap = map[string]string{
				"AppId":      AppId,
				"FolderPath": folderPath,
			}
		)

		bs := make([]byte, 0, 10240)
		buffer := bytes.NewBuffer(bs)

		os.MkdirAll(folderPath+"/app", os.ModePerm)
		f, _ := os.Create(fmt.Sprintf("%s/app/app.go", folderPath))
		err = appTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/cmd", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/cmd/server.go", folderPath))
		err = cmdTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/conf", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/conf/config_yaml.yaml", folderPath))
		f.Write([]byte(fmt.Sprintf(configYamlStr, AppId, Port)))

		os.MkdirAll(folderPath+"/doc", os.ModePerm)

		os.MkdirAll(folderPath+"/domain/demo", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/domain/demo/entity.go", folderPath))
		err = domainEntityTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/domain/demo/repo.go", folderPath))
		err = domainRepoTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/domain/demo/service.go", folderPath))
		err = domainServiceTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/handler", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/handler/handler.go", folderPath))
		err = handlerHandlerTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/handler/handler_test.go", folderPath))
		err = handlerTestTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/handler/hook.go", folderPath))
		err = handlerHookTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/handler/resource.go", folderPath))
		err = handlerResourceTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/infra/config", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/infra/config/config.go", folderPath))
		err = infraConfigTemplate.Execute(buffer, paramMap)
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/infra/err", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/infra/err/err.go", folderPath))
		err = infraErrTemplate.Execute(buffer, map[string]string{
			"AppId":      AppId,
			"FolderPath": folderPath,
			"Port":       fmt.Sprintf("%d", Port),
		})
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/infra/thirdparty/user", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/infra/thirdparty/user/user.go", folderPath))
		err = infraThirdPartyTemplate.Execute(buffer, map[string]string{
			"AppId":      AppId,
			"FolderPath": folderPath,
			"Port":       fmt.Sprintf("%d", Port),
		})
		utils.MustNil(err)
		utils.FormatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		os.MkdirAll(folderPath+"/proto", os.ModePerm)
		f, _ = os.Create(fmt.Sprintf("%s/proto/%s.proto", folderPath, AppId))
		err = protoTemplate.Execute(buffer, map[string]string{
			"AppId":      AppId,
			"FolderPath": folderPath,
		})
		utils.MustNil(err)
		f.Write(buffer.Bytes())
		//formatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		f, _ = os.Create(fmt.Sprintf("%s/Makefile", folderPath))
		err = makefileTemplate.Execute(buffer, map[string]string{
			"AppId":      AppId,
			"FolderPath": folderPath,
			"Port":       fmt.Sprintf("%d", Port),
			"App1Id":     fmt.Sprintf("%s-%s", serviceArr[0], serviceArr[1]),
		})
		utils.MustNil(err)
		f.Write(buffer.Bytes())
		//formatAndWrite(f, buffer.Bytes())
		buffer.Reset()

		err = os.Chdir(folderPath)
		utils.MustNil(err)
		c := exec.Command("make", "proto")
		err = c.Run()
		utils.MustNil(err)
	},
}
