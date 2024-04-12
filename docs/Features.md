# Features for Arknights Telegram Bot

This document provides an overview of the features available in the Arknights Telegram Bot.

## Private Chat Commands
```
/bind - Bind your Telegram account with your Arknights account.
/unbind - Unbind your Telegram account from your Arknights account.
/resume - Set your introduction for the group.
/reset_token - Reset your token status.
/btoken - Set your bilibili token.
/import_gacha - Import your gacha history.
/export_gacha - Export your gacha history.
/cancel - Cancel the current operation.
```
## All Chat Commands
```
/help - Display the help message.
/ping - Check the bot's status.
/sign - Daily check-in for Skland.
    /sign auto - Automatically check-in for Skland everyday.
    /sign stop - Cancel the automatic check-in.

/state - Show current statues of the binded Arknights account.
/box - Show the cards inventory for binded Arknights account.
    /box all - Show all rare level cards in the inventory.
    /box 6 - Show all 6-star cards in the inventory.
    /box 1,2 - Show all 1-star and 2-star cards in the inventory, note that the number of stars should be separated by commas.
/missing - Show the missing cards for binded Arknights account.
    /missing all - Show all missing rare level cards.
    /missing 6 - Show all missing 6-star cards.
    /missing 1,2 - Show all missing 1-star and 2-star cards, note that the number of stars should be separated by commas.
/card - Show the player card.
/base - Show the base status.
/gacha - Show the gacha history
/operator - Show the operator data, without any arguments, it will show the operator select helper. 
    /operator OperatorName - Show the operator data of the specified operator.
/skin - Show the skin data, without any arguments, it will show the skin select helper.
    /skin SkinName - Show the skin data of the specified skin.
/enemy - Show the enemy data, without any arguments, it will show the enemy select helper.
    /enemy EnemyName - Show the enemy data of the specified enemy.
/material - Show the material data, without any arguments, it will show the material select helper.
    /material MaterialName - Show the material data of the specified material.
/report - Report the Telegram account to the group administrator.
/quiz - Start a quiz game, player will given a random operator and need to choose the operator's name.
    /quiz h - Start a hard quiz game, player will given a random operator in pure dark mode and need to choose the operator's name.
    /quiz ex - Start a Special provided quiz game, player will given lower half of the random operator and need to choose the operator's name.
/redeem [code] - Redeem the code.
/headhunt - Emulate current gacha pool.
/recruit - Check the provided image to find what operator is possible to recruit with provided tags.
```
## Admin Commands
```
/news - Enable/Disable sending news to the group.
/quiz [start/stop] - Enable/Disable the quiz game.
/headhunt [start/stop] - Enable/Disable the headhunt game.
/reg - Set the message to be sent when a new user joins the group.
```
## Owner Commands
```
/update - Update the bot metadata.
/clear [key]- delete redis key
/kill - kill bot.
```