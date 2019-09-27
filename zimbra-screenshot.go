package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"strings"
)

func takeScreenshot(ConnectionSettings ConnectionServerConfig, v string){
	ConnectionSettings.adminAuthToken = getLoginToken("admin", ConnectionSettings)
	delegateToken := delegateAuthRequest(ConnectionSettings, v)

	zimbraMailHost := strings.Split(GetInfoRequest(ConnectionSettings, v, "host"), "home/")[0]

	ctx := context.Background()
	options := []chromedp.ExecAllocatorOption{
		chromedp.ProxyServer("socks5://"+ConnectionSettings.socksServerString),
	}
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	c, cc := chromedp.NewExecAllocator(ctx, options...)
	defer cc()
	// create context
	ctx, cancel := chromedp.NewContext(c)
	defer cancel()
	//var res string
	var buf []byte

	if err := chromedp.Run(ctx, fullScreenshot(zimbraMailHost+"/mail?adminPreAuth=1", 90, "ZM_AUTH_TOKEN", delegateToken, v, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(v+".png", buf, 0644); err != nil {
		log.Fatal(err)
	}else{
		log.Info("Saving screenshot for: "+v+ " to file "+v+".png")
	}
}

func fullScreenshot(urlstr string, quality int64, cookieName string, cookieValue string, email string, res *[]byte) chromedp.Tasks {
	cookieDomain := strings.Split(urlstr, "/")[2]
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Info("Setting cookie for user: "+email+ " on mailbox host: "+cookieDomain)
			success, err := network.SetCookie(cookieName, cookieValue).
				WithDomain(cookieDomain).
				WithHTTPOnly(true).
				Do(ctx)
			if err != nil {
				return err
			}
			if !success {
				return fmt.Errorf("could not set cookie %q to %q", cookieName, cookieValue)
			}

			return nil
		}),
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Info("Taking screenshot for user: "+email+" on host: "+cookieDomain)
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
