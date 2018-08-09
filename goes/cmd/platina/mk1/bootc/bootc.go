// Copyright © 2015-2018 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package bootc

import (
	"fmt"
	"os"

	"github.com/platinasystems/go/goes"
	"github.com/platinasystems/go/goes/lang"
)

func New() *Command { return new(Command) }

type Command struct {
	// Machines may use Hook to run something before redisd and other
	// daemons.
	Hook func() error

	// Machines may use ConfHook to run something after all daemons start
	// and before source of start command script.
	ConfHook func() error

	// GPIO init hook for machines than need it
	ConfGpioHook func() error

	g *goes.Goes
}

func (Command) String() string { return "bootc" }

//TODO CLEANUP
func (Command) Usage() string {
	return `
	bootc [register] [bootc] [dumpvars] [dashboard] [initcfg] [wipe]
	[getnumclients] [getclientdata] [getscript] [getbinary] [testscript]
	[test404] [dashboard9] [setsda6] [clrsda6] [setinstall] [clrinstall]
	[setsda1] [clrsda1] [readcfg] [setip] [setnetmask] [setgateway]
	[setkernel] [setinitrd] [getclientbootdata] [setpost] [vers]
	[checkfiles] [getfiles] [setdisable] [clrdisable] [pccinitfile]
	[setpccenb] [clrpccenb] [setpccip] [setpccport] [setpccsn]`
}

func (Command) Apropos() lang.Alt {
	return lang.Alt{
		lang.EnUS: "bootc provides wipe and access to bootc.cfg.",
	}
}

func (Command) Man() lang.Alt {
	return lang.Alt{
		lang.EnUS: `
description
	bootc provides wipe and access to bootc.cfg.`,
	}
}

func (c *Command) Goes(g *goes.Goes) { c.g = g }

func (c *Command) Main(args ...string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("args: missing")
	}
	mip := getMasterIP()

	//TODO Clean up these cases
	switch args[0] {
	case "register": //TODO CUT THIS
		mac := getMAC()
		ip := getIP()
		if _, _, err = register(mip, mac, ip); err != nil {
			return err
		}
	case "bootc":
		c.bootc()
		return nil
	case "vers":
		fmt.Println(verNum)
		return nil
	case "checkfiles":
		if err := checkFiles(); err != nil {
			return nil
		}
	case "getfiles":
		if err := getFiles(); err != nil {
			return nil
		}
	case "dumpvars": //TODO CUT THIS
		if err = dumpvars(mip); err != nil {
			return err
		}
	case "dashboard":
		if err = dashboard(mip); err != nil {
			return err
		}
	case "initcfg":
		if err = initCfg(); err != nil {
			return err
		}
	case "getnumclients":
		if err = getnumclients(mip); err != nil {
			return err
		}
	case "getclientdata":
		if err = getclientdata(mip, 3); err != nil {
			return err
		}
	case "getclientbootdata":
		if err = getclientbootdata(mip, 3); err != nil {
			return err
		}
	case "getscript":
		if err = getscript(mip, "testscript"); err != nil {
			return err
		}
	case "getbinary":
		if err = getbinary(mip, "test.bin"); err != nil {
			return err
		}
	case "testscript":
		if err = runScript("testscript"); err != nil {
			return err
		}
	case "test404":
		if err = test404(mip); err != nil {
			return err
		}
	case "dashboard9":
		if err = dashboard("192.168.101.129"); err != nil {
			return err
		}
	case "setsda6":
		if err = setSda6Count(args[1]); err != nil {
			return err
		}
	case "clrsda6":
		if err = clrSda6Count(); err != nil {
			return err
		}
	case "setinstall":
		if err = setInstall(); err != nil {
			return err
		}
	case "clrinstall":
		if err = clrInstall(); err != nil {
			return err
		}
	case "setsda1":
		if err = setSda1Flag(); err != nil {
			return err
		}
	case "clrsda1":
		if err = clrSda1Flag(); err != nil {
			return err
		}
	case "setdisable":
		if err = setDisable(); err != nil {
			return err
		}
	case "clrdisable":
		if err = clrDisable(); err != nil {
			return err
		}
	case "setpccenb":
		if err = setPccEnb(); err != nil {
			return err
		}
	case "clrpccenb":
		if err = clrPccEnb(); err != nil {
			return err
		}
	case "pccip":
		if err = setPccIP(args[1]); err != nil {
			return err
		}
	case "pccport":
		if err = setPccPort(args[1]); err != nil {
			return err
		}
	case "pccsn":
		if err = setPccSN(args[1]); err != nil {
			return err
		}
	case "readcfg":
		if err := readCfg(); err != nil {
			return err
		}
		fmt.Printf("%+v\n", Cfg)
	case "setip":
		if err = setIp(args[1]); err != nil {
			return err
		}
	case "setnetmask":
		if err = setNetmask(args[1]); err != nil {
			return err
		}
	case "setgateway":
		if err = setGateway(args[1]); err != nil {
			return err
		}
	case "setkernel":
		if err = SetSda6K(args[1]); err != nil {
			return err
		}
	case "setinitrd":
		if err = SetSda6I(args[1]); err != nil {
			return err
		}
	case "setpost":
		if err = setPostInstall(); err != nil {
			return err
		}
	case "wipe":
		if len(os.Args) >= 3 && args[1] == "sda6" {
			if err = Wipe(); err != nil {
				return err
			}
			reboot()
		}
		fmt.Println("Type: 'bootc wipe sda6' to re-install debian on sda6")
		return nil

		//commands to test pcc messages
	case "pccinitfile":
		err := pccInitFile()
		if err != nil {
			return err
		}
		return nil
	case "pccinit":
		err := pccInit()
		if err != nil {
			return err
		}
		return nil
	case "pcc1a":
		data, err := doPost(BSTAT, "goes-boot.booting")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	case "pcc1b":
		data, err := doPost(BSTAT, "goes-boot.operational")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	case "pcc1c":
		data, err := doPost(BSTAT, "goes-boot.wiping")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	case "pcc2":
		data, err := doPost(RDCFG, "")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	case "pcc3":
		data, err := doPost(REGIS, "")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	case "pcc4":
		data, err := doPost(UNREG, "")
		if err != nil {
			return err
		}
		fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
		return nil
	default:
	}
	return nil
}
