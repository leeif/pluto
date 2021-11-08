## Sign in with WeChat QRCode

https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html

1. Pluto 收到微信请求，带有参数 code 和 state
2. state 是 base64 编码的 JSON 字符串，
  - app: 对应 pluto app 的 name, e.g. org.mushare.easyjapanese
  - redirect_url: 跳转 URL，pluto 登录成功以后跳转的 URL
3. 根据 code 进行微信登陆
4. 302 跳转到 redirectURL + "?token=...."
