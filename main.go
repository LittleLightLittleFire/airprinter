package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/OpenPrinting/goipp"
)

func main() {
	var ipArg string
	var port int
	var instance string
	var newHost string

	var proxyHost string
	flag.StringVar(&proxyHost, "printerHost", "http://192.168.1.77:631", "real host of the printer")

	flag.StringVar(&ipArg, "newIP", "192.168.1.2", "new ip of the printer")
	flag.IntVar(&port, "newPort", 631, "new port of the printer")
	flag.StringVar(&newHost, "newHost", "A6PRINTER", "new host of the printer")

	flag.StringVar(&instance, "instance", "a6printer", "name of the printer")
	flag.Parse()

	var txt []string
	txt = append(txt, "txtvers=1")
	txt = append(txt, "qtotal=1")
	txt = append(txt, "pdl=application/octet-stream,image/urf,image/jpeg,image/pwg-raster")
	txt = append(txt, "rp=ipp/print")
	txt = append(txt, "note=")
	txt = append(txt, "ty=Brother HL-L3230CDW series")
	txt = append(txt, "product=(Brother HL-L3230CDW series)")
	txt = append(txt, fmt.Sprintf("adminurl=http://%v.local.:%v/net/net/airprint.html", newHost, port))
	txt = append(txt, "priority=25")
	txt = append(txt, "usb_MFG=Brother")
	txt = append(txt, "usb_MDL=HL-L3230CDW series")
	txt = append(txt, "usb_CMD=PJL,PCL,PCLXL,URF")
	txt = append(txt, "Color=T")
	txt = append(txt, "Copies=T")
	txt = append(txt, "Duplex=T")
	txt = append(txt, "Fax=F")
	txt = append(txt, "Scan=F")
	txt = append(txt, "PaperCustom=T")
	txt = append(txt, "Binary=T")
	txt = append(txt, "Transparent=T")
	txt = append(txt, "TBCP=F")
	txt = append(txt, "URF=SRGB24,W8,CP1,IS4-1,MT1-3-4-5-8-11,OB10,PQ4,RS600,V1.4,DM1")
	txt = append(txt, "kind=document,envelope,label,postcard")
	txt = append(txt, "PaperMax=legal-A4")
	txt = append(txt, "UUID=dbbd227c-d891-4474-9c92-7f20dcb234f4")
	txt = append(txt, "print_wfds=T")
	txt = append(txt, "mopria-certified=1.3")

	var ips []net.IP
	ip := net.ParseIP(ipArg)
	ips = append(ips, ip)

	service, err := NewMDNSService(instance, "_ipp._tcp", "", newHost+".local.", port, ips, txt)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the mDNS server, defer shutdown
	server, _ := NewServer(&Config{
		Zone: service,
	})
	defer server.Shutdown()

	// Create the IPP proxy
	proxyHostURL, err := url.Parse(proxyHost)
	if err != nil {
		log.Fatalln(err)
	}

	rp := httputil.NewSingleHostReverseProxy(proxyHostURL)
	rp.ModifyResponse = func(resp *http.Response) error {
		if resp.Header.Get("Content-Type") != "application/ipp" {
			return nil
		}

		var msg goipp.Message
		if err := msg.Decode(resp.Body); err != nil {
			return err
		}

		for i, group := range msg.Groups {
			for j, attr := range group.Attrs {
				if attr.Name == "media-col-ready" {
					for k, value := range attr.Values {
						if value.T != goipp.TagBeginCollection {
							continue
						}

						log.Println("attr", i, j, k, attr.Name, attr.Values)

						var mediaColReady goipp.Collection
						for _, v := range value.V.(goipp.Collection) {
							switch v.Name {
							case "media-size":
								var mediaSize goipp.Collection
								mediaSize.Add(goipp.MakeAttribute("x-dimension", goipp.TagInteger, goipp.Integer(10500)))
								mediaSize.Add(goipp.MakeAttribute("y-dimension", goipp.TagInteger, goipp.Integer(14800)))
								mediaColReady.Add(goipp.MakeAttribute(v.Name, goipp.TagBeginCollection, mediaSize))

							case "media-source":
								mediaColReady.Add(goipp.MakeAttribute(v.Name, goipp.TagKeyword, goipp.String("manual")))
							default:
								mediaColReady.Add(v)
							}
						}
						msg.Groups[i].Attrs[j].Values[k].V = mediaColReady
					}
				}
			}
		}

		log.Println("Resp:", msg.Groups)

		data, err := msg.EncodeBytes()
		if err != nil {
			return err
		}

		resp.Body = io.NopCloser(bytes.NewReader(data))
		resp.Header.Set("Content-Length", fmt.Sprint(len(data)))
		return nil
	}

	host := fmt.Sprintf("%v:%v", ipArg, port)
	host80 := fmt.Sprintf("%v:%v", ipArg, 80)
	log.Println("Reverse proxy:", host, "=>", proxyHost)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[http]", r.RequestURI)
		rp.ServeHTTP(w, r)
	})
	go http.ListenAndServe(host, nil)
	go http.ListenAndServe(host80, nil)

	select {}
}
