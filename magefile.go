// +build mage

package main

import "github.com/magefile/mage/sh"

func Base() error {
	err := sh.Run("packr2", "clean")
	if err != nil {
		return err
	}
	err = sh.Run("make", "-C", "proto", "proto")
	if err != nil {
		return err
	}
	return nil
}
func Build() error {
	err := Base()
	if err != nil {
		return err
	}
	err = sh.Run("packr2")
	if err != nil {
		return err
	}
	err = sh.Run("go", "install", ".")
	if err != nil {
		sh.Run("packr2", "clean")
		return err
	}
	err = sh.Run("packr2", "clean")
	if err != nil {
		return err
	}
	return nil
}
func Proto() error {
	err := Base()
	if err != nil {
		return err
	}
	err = sh.Run("make", "-C", "example", "proto")
	if err != nil {
		return err
	}
	return nil
}
func Example() error {
	err := Build()
	if err != nil {
		return err
	}
	err = Proto()
	if err != nil {
		return err
	}
	err = sh.Run("packr2", "clean")
	if err != nil {
		return err
	}
	return nil
}
func Clean() error {
	err := sh.Run("make", "-C", "example", "clean")
	if err != nil {
		return err
	}
	err = sh.Run("packr2", "clean")
	if err != nil {
		return err
	}
	return nil
}
