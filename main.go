package main

import (
	"os/exec"
	"strings"

	"github.com/yutaro/slack-cmd-go"
)

var (
	templates = map[string]string{
		"text":    "text.tpe",
		"text-qr": "text-qr.tpe",
	}
)

var exepath = "C:/Program Files (x86)/KING JIM/TEPRA SPC10/SPC10.exe"

func main() {
	conf := scmd.LoadToml("config.toml")
	bot := scmd.New(conf.TOKEN)
	tepra := bot.NewCmdGroup("tepra")

	imgpath := imgPath("result")

	tepra.Cmd("print", []string{"print message"},
		func(c *scmd.Context) {
			args := c.GetArgs()
			options := c.GetOptions()
			flags := c.GetFlags()
			mes := strings.Join(args, " ")

			tpe := "text"

			prints := []string{mes}

			if qr, ok := options["qr"]; ok {
				url := urlConv(qr)
				prints = append(prints, url)
				tpe = "text_qr"
			}

			csvpath := writeCsv(prints)
			tpepath := tpePath(tpe)

			n := "1"
			if num, ok := options["n"]; ok {
				n = num
			}

			cmd := exec.Command(exepath, "-p", tpepath+","+csvpath+","+n, "-B", "-a "+imgpath)
			cmd.Run()
			//cmd.Wait()

			if flags["t"] {
				c.SendFile(imgpath)
			}

			c.SendMessage(mes)
		})

	tepra.Cmd("qrcode", []string{"print only qrcode"},
		func(c *scmd.Context) {
			args := c.GetArgs()
			if len(args) < 1 {
				return
			}

			qr := args[0]
			url := urlConv(qr)

			prints := []string{url}

			csvpath := writeCsv(prints)
			tpepath := tpePath("qr")

			n := "1"
			if num, ok := c.GetOptions()["n"]; ok {
				n = num
			}

			cmd := exec.Command(exepath, "-p", tpepath+","+csvpath+","+n)
			cmd.Run()

			c.SendMessage(url)
		})

	tepra.Cmd("image", []string{"print image", "don't work now..."},
		func(c *scmd.Context) {
			c.SendMessage("preparing now...")
		})

	tepra.Cmd("status", []string{"check status"},
		func(c *scmd.Context) {
			c.SendMessage("ok")
		})

	bot.Start()
}
