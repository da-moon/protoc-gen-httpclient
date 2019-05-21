package main

import (
	"encoding/json"
	"log"
	"strings"

	pb "github.com/bifrostcloud/protoc-gen-httpclient/proto"
	"github.com/fatih/camelcase"
	proto "github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/palantir/stacktrace"
)

type base struct {
	Package  string
	Imports  []pkg
	Services []service
}

type service struct {
	UpperCamelCaseServiceName string
	LowerCamelCaseServiceName string
	Auth                      string
	Methods                   []rpc
}
type fieldImport struct {
	Name string
	Tag  string
}
type fields struct {
	// Field Structs
	FieldImport []fieldImport
	// Proto Message Name
	Type string

	UpperCamelCase string
	// Keeps the fields name as a lowercase of message field in proto file
	Base []string
	// turns the field name into a single lowercase string without any seperators
	Lowercase []string
	// turns the field name into a single lowercase with dot seperators
	DotNotation []string
	// turns the field name into a lowercase, paramcase string
	ParamCase []string
}
type rpc struct {
	UpperCamelCaseMethodName  string
	UpperCamelCaseServiceName string
	InputType                 string
	InputFields               fields
	OutputType                string
	Auth                      string
	RequestOptions            pb.RequestOptions
}
type Input struct {
	FieldName  string
	FieldValue string
}
type pkg struct {
	PackageName string
	PackagePath string
}

func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range targets {
		name := p.Context.OutputPath(file).SetExt(".httpclient.go").String()
		b := base{Package: p.Context.PackageName(file).String()}
		imports := map[string]pkg{
			"net": pkg{
				PackagePath: "net/http",
			},
			"stacktrace": pkg{
				PackageName: "stacktrace",
				PackagePath: "github.com/palantir/stacktrace",
			},
		}
		for _, srv := range file.Services() {
			opt := srv.Descriptor().GetOptions()
			option, err := proto.GetExtension(opt, pb.E_ServiceOptions)
			if err != nil {
				if err == proto.ErrMissingExtension {
					continue
				}
				// log.Fatal(stacktrace.NewError(err.Error()))
			}
			byteData, err := json.Marshal(option)
			if err != nil {
				log.Fatal(stacktrace.NewError(err.Error()))
			}
			srvOpts := pb.ServiceOptions{}
			err = json.Unmarshal(byteData, &srvOpts)
			if err != nil {
				log.Fatal(stacktrace.NewError(err.Error()))
			}
			s := service{}
			s.UpperCamelCaseServiceName = srv.Name().UpperCamelCase().String()
			s.LowerCamelCaseServiceName = srv.Name().LowerCamelCase().String()
			s.Auth = strings.ToLower(srvOpts.Auth)
			for _, method := range srv.Methods() {
				opt := method.Descriptor().GetOptions()
				option, err := proto.GetExtension(opt, pb.E_RequestOptions)
				if err != nil {
					if err == proto.ErrMissingExtension {
						continue
					}
					// log.Fatal(stacktrace.NewError(err.Error()))
				}
				byteData, err := json.Marshal(option)
				if err != nil {
					log.Fatal(stacktrace.NewError(err.Error()))
				}
				clientOpts := pb.RequestOptions{}
				err = json.Unmarshal(byteData, &clientOpts)
				if err != nil {
					log.Fatal(stacktrace.NewError(err.Error()))
				}
				clientOpts.Target = strings.ToLower(srvOpts.Endpoint + clientOpts.Target)
				clientOpts.ClientType = strings.ToLower(clientOpts.ClientType)
				clientOpts.Method = strings.ToUpper(clientOpts.Method)
				if clientOpts.Method == "POST" || clientOpts.Method == "PUT" {
					imports["io"] = pkg{
						PackagePath: "io",
					}
					imports["ioutil"] = pkg{
						PackagePath: "io/ioutil",
					}
					imports["url"] = pkg{
						PackagePath: "net/url",
					}
					imports["strings"] = pkg{
						PackagePath: "strings",
					}
				}

				imports["utils"] = pkg{
					PackageName: "utils",
					PackagePath: "github.com/bifrostcloud/protoc-gen-httpclient/pkg/utils",
				}
				imports["json"] = pkg{
					PackagePath: "encoding/json",
				}

				if clientOpts.ClientType == "circuit-breaker" {
					imports["circuit-breaker"] = pkg{
						PackageName: "cb",
						PackagePath: "github.com/bifrostcloud/protoc-gen-httpclient/pkg/client/circuit-breaker",
					}
				} else if clientOpts.ClientType == "basic" {
					imports["basic"] = pkg{
						PackageName: "basic",
						PackagePath: "github.com/bifrostcloud/protoc-gen-httpclient/pkg/client/basic",
					}
				}
				ms := p.Context.Name(method).UpperCamelCase().String()
				splitted := camelcase.Split(ms)
				firstElem := strings.ToLower(splitted[0])
				if firstElem == "put" || firstElem == "post" || firstElem == "get" || firstElem == "delete" {
					splitted[0] = ""
				}
				lastElem := strings.ToLower(splitted[len(splitted)-1])
				if lastElem == "put" || lastElem == "post" || lastElem == "get" || lastElem == "delete" {
					splitted[len(splitted)-1] = ""
				}
				ms = strings.Join(splitted, "")
				upperCamelCaseMethodName := pgs.Name(strings.ToLower(clientOpts.Method)).UpperCamelCase().String() + ms
				p.Context.Name(method.Input()).String()

				r := rpc{}
				r.UpperCamelCaseServiceName = srv.Name().UpperCamelCase().String()
				r.UpperCamelCaseMethodName = upperCamelCaseMethodName

				r.InputType = p.Context.Name(method.Input()).String()
				if !method.Input().BuildTarget() {
					path := p.Context.ImportPath(method.Input()).String()
					imports[path] = pkg{
						PackagePath: path,
					}

					r.InputType = p.Context.PackageName(method.Input()).String() + "." + p.Context.Name(method.Input()).String()
				}
				r.populateInputFields(method.Input().Fields())
				r.OutputType = p.Context.Name(method.Output()).String()
				r.Auth = strings.ToLower(srvOpts.Auth)
				r.RequestOptions = clientOpts
				if !method.Output().BuildTarget() {
					path := p.Context.ImportPath(method.Output()).String()
					imports[path] = pkg{
						PackagePath: path,
					}

					r.OutputType = p.Context.PackageName(method.Output()).String() + "." + p.Context.Name(method.Output()).String()
				}

				s.Methods = append(s.Methods, r)
			}
			b.Services = append(b.Services, s)
		}

		if len(b.Services) == 0 {
			continue
		}

		for _, pkg := range imports {
			b.Imports = append(b.Imports, pkg)
		}

		p.OverwriteGeneratorTemplateFile(
			name,
			template.Lookup("Base"),
			&b,
		)
	}

	return p.Artifacts()
}
func (r *rpc) populateInputFields(inputFields []pgs.Field) {
	for _, field := range inputFields {
		r.InputFields.Type = r.InputType
		r.InputFields.FieldImport = append(r.InputFields.FieldImport, fieldImport{
			Name: field.Name().UpperCamelCase().String(),
			Tag:  field.Name().String(),
		})
		r.InputFields.Base = append(r.InputFields.Base, strings.ToLower(field.Name().String()))
		r.InputFields.Lowercase = append(r.InputFields.Lowercase, strings.ToLower(field.Name().LowerCamelCase().String()))
		r.InputFields.DotNotation = append(r.InputFields.DotNotation, strings.ToLower(field.Name().LowerDotNotation().String()))
		spl := camelcase.Split(field.Name().UpperCamelCase().String())
		r.InputFields.ParamCase = append(r.InputFields.ParamCase, strings.ToLower(strings.Join(spl, "-")))
	}
}

// md, _ := example.EncodeVersionRequestToMetadata(example.VersionRequest{
// 	Status:  "some value status",
// 	Boolean: true,
// 	Integer: 328,
// })
// result := map[string]string{}
// for k, v := range md {
// 	result[k] = fmt.Sprintf("%v", v)
// }
// log.Printf("%#v", result)
