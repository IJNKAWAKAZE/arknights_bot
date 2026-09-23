package media

import "errors"

// Close 在截图任务和 HTTP 请求全部结束后调用，释放浏览器及驱动资源。
func Close() error {
	browserLock.Lock()
	defer browserLock.Unlock()
	var errs []error
	for scale, ctx := range contexts {
		if !ctx.IsClosed() {
			errs = append(errs, ctx.Close())
		}
		delete(contexts, scale)
	}
	if browser != nil {
		errs = append(errs, browser.Close())
		browser = nil
	}
	if playwrightDriver != nil {
		errs = append(errs, playwrightDriver.Stop())
		playwrightDriver = nil
	}
	return errors.Join(errs...)
}
