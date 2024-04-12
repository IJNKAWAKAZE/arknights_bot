# 使用说明

现在能用的功能

## 私聊指令
```
/bind - 绑定森空岛账号(支持多账号与单账号多角色)
/unbind - 解绑角色
/resume - 设置角色签名
/reset_token - 重设森空岛Token
/btoken - 设置BToken(B服用户才需要)
/import_gacha - 导入抽卡记录
/export_gacha - 导出抽卡记录
/cancel - 取消当前操作
```
## 普通指令
```
/help - 查看使用说明
/ping - 存活测试
/sign - 森空岛签到
    /sign auto - 开启自动签到(自动签到包含该tg号下绑定的所有角色)
    /sign stop - 关闭自动签到

/state - 查看游戏账号当前状态(森空岛小组件上那些)
/box - 查看游戏干员列表(默认只查看6星)
    /box all - 查看所有干员
    /box 6 - 查看对应星级干员
    /box 1,2 - 查看多个星级干员
/missing - 查看游戏未获取干员(默认只查看6星)
    /missing all - 查看所有未获取干员
    /missing 6 - 查看对应星级未获取干员
    /missing 1,2 - 查看多个星级未获取干员
/card - 查看角色名片(私聊设置的签名会在这里显示)
/base - 查看基建状态
/gacha - 查看抽卡分析
/operator - 查看干员信息(不输入名称显示干员搜索按钮)
    /operator 阿米娅 - 查看干员信息
/skin - 查看干员皮肤(不输入名称显示干员搜索按钮)
    /skin 阿米娅 - 查看干员皮肤
/enemy - 查看敌人信息(不输入名称显示敌人搜索按钮)
    /enemy 源石虫 - 查看敌人信息
/material - 查看材料刷取推荐(不输入名称显示材料搜索按钮)
    /material 装置 - 查看材料刷取推荐
/report - 回复群组消息进行举报(会自动@所有管理员)
/quiz - 猜干员小游戏(默认立绘模式)
    /quiz h - 黑白剪影模式
    /quiz ex - 看脚猜干员模式(合乎粥礼)
/redeem [code] - 兑换CDK(如果使用的人过多会触发风控，需要稍后才能恢复)
/headhunt - 寻访模拟(需要自己更新卡池信息)
/recruit - 公招计算(指令需附带图片一起发送)
    /recruit jp 日服公招计算
```
## 管理员指令
```
/news - 开启/关闭B站动态推送(默认关闭状态)
/quiz [start/stop] - 开启/关闭猜干员小游戏
/headhunt [start/stop] - 开启/关闭寻访模拟
/reg - 回复一条消息设置为群规
```
## 机器人拥有者的指令
```
/update - 更新数据源(首次运行项目需要先执行一次，后续每周五凌晨会自动更新)
/clear [key]- 根据key删除redis缓存
/kill - 杀死机器人
```