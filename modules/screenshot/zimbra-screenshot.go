package screenshot

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
	"../common"
	"../zimbra/soap"
)
var Log logrus.Logger
func TakeScreenshot(ConnectionSettings common.ConnectionServerConfig, v string){
	ConnectionSettings.AdminAuthToken = soap.GetLoginToken("admin", ConnectionSettings)
	delegateToken := soap.DelegateAuthRequest(ConnectionSettings, v)
	infoRequest := soap.GetInfoRequest(ConnectionSettings, v, "host")
	zimbraMailHost := strings.Split(infoRequest, "home/")[0]
	ctx := context.Background()
	if ConnectionSettings.UseSocks5Proxy ==true{
		options := []chromedp.ExecAllocatorOption{
			chromedp.ProxyServer("socks5://"+ConnectionSettings.SocksServerString),
		}
		options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
		c, cc := chromedp.NewExecAllocator(ctx, options...)
		defer cc()
		ctx, cancel := chromedp.NewContext(c)
		defer cancel()
		var buf []byte
		if err := chromedp.Run(ctx, fullScreenshot(zimbraMailHost+"/mail?adminPreAuth=1", 90, "ZM_AUTH_TOKEN", delegateToken, v, &buf)); err != nil {
			cancel()
			cc()
			return
		}
		if err := ioutil.WriteFile(v+".png", buf, 0644); err != nil {
			for {
				if err = ioutil.WriteFile(v+".png", buf, 0644); err != nil {
						Log.Error("[Screenshotter] Connection Issue, throttling...")
						time.Sleep(2 * time.Second)
				}else{
					break
				}
			}

		}else{
			Log.Info("["+v+"] saving screenshot to file "+v+".png")
		}
	}else{
		var options []chromedp.ExecAllocatorOption
		options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
		c, cc := chromedp.NewExecAllocator(ctx, options...)
		defer cc()
		ctx, cancel := chromedp.NewContext(c)
		defer cancel()
		var buf []byte

		if err := chromedp.Run(ctx, fullScreenshot(zimbraMailHost+"/mail?adminPreAuth=1", 90, "ZM_AUTH_TOKEN", delegateToken, v, &buf)); err != nil {
			Log.Fatal(err)
		}
		//endTime = time.Now()
		if err := ioutil.WriteFile(v+".png", buf, 0644); err != nil {
			Log.Fatal(err)
		}else{
			Log.Info("Saving screenshot for: "+v+ " to file "+v+".png")
		}
	}
	host := strings.Split(zimbraMailHost, ":")[1]
	host = strings.Replace(host, "//","",-1)
}

func fullScreenshot(urlstr string, quality int64, cookieName string, cookieValue string, email string, res *[]byte) chromedp.Tasks {
	cookieDomain := strings.Split(urlstr, "/")[2]
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			Log.Info("["+email+"] Setting cookie on mailbox host: "+cookieDomain)
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
