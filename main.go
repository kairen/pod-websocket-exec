package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func markRequiredFlags(opts *ExecOptions) {
	if len(opts.Pod) == 0 || opts.Command == nil {
		fmt.Fprintf(os.Stderr, "Fatal error missing some flags:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func parseFlags(opts *ExecOptions) {
	flag.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "Absolute path to the kubeconfig file.")
	flag.StringVarP(&opts.Namespace, "namespace", "n", "default", "Exec pod in which namespace.")
	flag.StringVarP(&opts.Pod, "pod", "p", "", "Exec pod name.（required）")
	flag.StringVarP(&opts.Container, "container", "c", "", "Exec container name in a pod.")
	flag.BoolVarP(&opts.TTY, "tty", "t", false, "Stdin is a TTY.")
	flag.BoolVarP(&opts.Stdin, "stdin", "i", false, "Pass stdin to the container.")
	flag.StringSliceVar(&opts.Command, "command", nil, "Exec commands.（required）")
	flag.Parse()
	markRequiredFlags(opts)
}

func replaceHomePath(path string) string {
	p := path
	if index := strings.Index(p, "~/"); index > -1 {
		home, _ := homedir.Dir()
		p = filepath.Join(home, p[index+len("~/"):])
	}
	return p
}

func main() {
	opts := &ExecOptions{}
	parseFlags(opts)

	kubeconfig = replaceHomePath(kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalln(err)
	}

	wrt, err := ExecRoundTripper(config, WebsocketCallback)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := ExecRequest(config, opts)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := wrt.RoundTrip(req); err != nil {
		log.Fatalln(err)
	}
}
