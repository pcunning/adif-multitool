// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"

	"github.com/flwyd/adif-multitool/adif/spec"
	"github.com/flwyd/adif-multitool/cmd"
)

type cmdConfig struct {
	cmd.Command
	Configure func(*cmd.Context, *flag.FlagSet)
}

var (
	catConf = cmdConfig{Command: cmd.Cat}

	editConf = cmdConfig{Command: cmd.Edit,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.EditContext{
				Add:    cmd.NewFieldAssignments(cmd.ValidateAlphanumName),
				Set:    cmd.NewFieldAssignments(cmd.ValidateAlphanumName),
				Remove: make(cmd.FieldList, 0)}
			// fs.Var(&cctx.If, "if", "Only edit records where `field=value` is already set (repeatable)")
			fs.Var(cctx.Cond.IfFlag(), "if", "Only edit records where `condition` is true (repeatable)")
			fs.Var(cctx.Cond.IfNotFlag(), "if-not", "Only edit records where `condition` is false (repeatable)")
			fs.Var(cctx.Cond.OrIfFlag(), "or-if", "Only edit records where `condition` is true or any previous --if group is true (repeatable)")
			fs.Var(cctx.Cond.OrIfNotFlag(), "or-if-not", "Only edit records where `condition` is false or any previous --if group is true (repeatable)")
			fs.Var(&cctx.Add, "add", "Add `field=value` if field is not already in a record (repeatable)")
			fs.Var(&cctx.Set, "set", "Set `field=value` for all records (repeatable)")
			fs.Var(&cctx.Remove, "remove", "Remove `fields` from all records (comma-separated, repeatable)")
			fs.BoolVar(&cctx.RemoveBlank, "remove-blank", false, "Remove all blank fields")
			fs.Var(&cctx.FromZone, "time-zone-from", "Adjust times and dates from this time `zone` into -time-zone-to (default UTC)")
			fs.Var(&cctx.ToZone, "time-zone-to", "Adjust times and dates into this time `zone` from -time-zone-from (default UTC)")
			ctx.CommandCtx = &cctx
		}}

	findConf = cmdConfig{Command: cmd.Find,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.FindContext{}
			fs.Var(cctx.Cond.IfFlag(), "if", "Include records where `condition` is true (repeatable)")
			fs.Var(cctx.Cond.IfNotFlag(), "if-not", "Include records where `condition` is false (repeatable)")
			fs.Var(cctx.Cond.OrIfFlag(), "or-if", "Include records where `condition` is true or any previous --if group is true (repeatable)")
			fs.Var(cctx.Cond.OrIfNotFlag(), "or-if-not", "Include records where `condition` is false or any previous --if group is true (repeatable)")
			ctx.CommandCtx = &cctx
		}}

	fixConf = cmdConfig{Command: cmd.Fix}

	helpConf = cmdConfig{Command: cmd.Command{
		Name: "help", Description: "Print program or command usage information",
		Run: func(*cmd.Context, []string) error {
			// handled specially by main
			return nil
		}}}

	inferConf = cmdConfig{Command: cmd.Infer,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.InferContext{}
			fs.Var(&cctx.Fields, "fields", "Comma-separated or multiple instance field `names` to infer if absent")
			fs.BoolVar(&cctx.CommentLog, "comment-log", false, "Add record comments with a list of successfully inferred fields")
			ctx.CommandCtx = &cctx
		}}

	saveConf = cmdConfig{Command: cmd.Save,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.SaveContext{}
			fs.BoolVar(&cctx.CreateDirectory, "create-dirs", false, "Create any needed parent directories of the output file(s)")
			fs.BoolVar(&cctx.Quiet, "quiet", false, "Do not print record counts and file names to stderr")
			fs.BoolVar(&cctx.OverwriteExisting, "overwrite-existing", false, "Overwrite output file if it already exists")
			fs.BoolVar(&cctx.WriteIfEmpty, "write-if-empty", false, "Write output file even if standard input has no records")
			ctx.CommandCtx = &cctx
		}}

	selectConf = cmdConfig{Command: cmd.Select,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.SelectContext{Fields: make(cmd.FieldList, 0, 16)}
			fs.Var(&cctx.Fields, "fields", "Comma-separated or multiple instance field `names` to include in output")
			ctx.CommandCtx = &cctx
		}}

	sortConf = cmdConfig{Command: cmd.Sort,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.SortContext{Fields: make(cmd.FieldList, 0, 16)}
			fs.Var(&cctx.Fields, "fields", "Comma-separated or multiple instance field `names` to sort by")
			ctx.CommandCtx = &cctx
		}}

	validateConf = cmdConfig{Command: cmd.Validate}

	versionConf = cmdConfig{Command: cmd.Command{
		Name: "version", Description: "Print program version information",
		Run: func(*cmd.Context, []string) error {
			fmt.Printf("%s version %s\n", programName, version)
			fmt.Printf("ADIF version %s from %s\n", spec.ADIFVersion, spec.ADIFSpecURL)
			fmt.Printf("See %s for details\n", helpUrl)
			return nil
		}}}

	cmds = []cmdConfig{
		catConf,
		editConf,
		findConf,
		fixConf,
		helpConf,
		inferConf,
		saveConf,
		selectConf,
		sortConf,
		validateConf,
		versionConf,
	}
)

func commandNamed(name string) (cmdConfig, bool) {
	for _, c := range cmds {
		if c.Name == name {
			return c, true
		}
	}
	return cmdConfig{}, false
}

func commandNames() []string {
	res := make([]string, len(cmds))
	for i, c := range cmds {
		res[i] = c.Name
	}
	return res
}
