package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {

	fmt.Println("ISO Installer")
	fmt.Println("=============")
	fmt.Println()

	disks, err := getDisks()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Found %d disk(s):\n", len(disks))
	for i, disk := range disks {
		fmt.Printf("  %d. %s\n", i+1, disk.Name)
	}

	if len(disks) == 0 {
		fmt.Println("No disks found!")
		return
	}

	var i int
ReadInput:
	fmt.Printf("\nEnter disk number from above (1-%d)? ", len(disks))
	_, err = fmt.Scanf("%d", &i)
	if err != nil {
		fmt.Println(err.Error())
		goto ReadInput
	}
	if i < 1 || i > len(disks) {
		fmt.Println("invalid value!")
		goto ReadInput
	}
	i--

	fmt.Printf("You choose %s\n", disks[i].Name)

	fmt.Printf("OVERWRITE ALL DATA ON DISK %s (y/N)? ", disks[i].Name)
	var r rune
	_, err = fmt.Scanf("%c", &r)
	if err != nil || r != 'y' || r == 'Y' {
		fmt.Println("Goodbye")
		return
	}

	var imgFile string
	imgFile = "/var/disk.img"

	// @TODO offer choice of img files

	args := []string{"dd", "if=" + imgFile, "of=/dev/" + disks[i].Name, "status=progress"}

	fmt.Println("Executing: ", strings.Join(args, " "))

	fmt.Println("@TODO...")

	// @TODO Resync

	// @TODO Get partitions

	// @TODO Offer to expand last partition

}

type Disk struct {
	Name       string
	Partitions []Partition
}

type Partition struct {
	Name       string
	MountPoint string
}

func getDisks() ([]Disk, error) {
	disks := []Disk{}

	b, err := exec.Command("lsblk", "-r").Output()
	if err != nil {
		return disks, err
	}

	/*
		NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINTS
		...
		loop30 7:30 0 88.1M 1 loop /snap/zulip/51
		nvme0n1 259:0 0 1.8T 0 disk
		nvme0n1p1 259:1 0 512M 0 part /boot/efi
		nvme0n1p2 259:2 0 1.8T 0 part /var/snap/firefox/common/host-hunspell\x0a/
	*/

	for i, line := range bytes.Split(b, []byte{0x0a}) {
		if len(line) == 0 {
			continue
		}
		e := strings.Split(string(line), " ")
		if len(e) < 6 {
			fmt.Printf("Line %d: expected at least 6 cols: '%s'\n", i, line)
			continue
		}
		if e[5] == "disk" {
			disks = append(disks, Disk{
				Name: e[0],
			})
		}
	}

	return disks, nil
}
