> **Warning**
> 
> Sea Dinner API is no longer available publicly, hence SeaHungerGamesBot is dead now :( 
> 
> R.I.P March ~ July 2022

<h1 align = "center"> SeaHungerGames </h1>
<p align="center"><img src = "static/banner.gif"></p>

<div align="center">

[![Deployment](https://github.com/aaronangxz/SeaDinner/actions/workflows/deploy.yaml/badge.svg?branch=master)](https://github.com/aaronangxz/SeaDinner/actions/workflows/deploy.yaml)
[![Test and coverage](https://github.com/aaronangxz/SeaDinner/actions/workflows/codecov.yaml/badge.svg)](https://github.com/aaronangxz/SeaDinner/actions/workflows/codecov.yaml)
[![codecov](https://codecov.io/gh/aaronangxz/SeaDinner/branch/master/graph/badge.svg?token=AR5L758FVV)](https://codecov.io/gh/aaronangxz/SeaDinner)

</div>

<div align="center"> <em>"May the odds be ever in your favor."</em> </div>

<h1> How To Use </h1>

1. Chat with `SeaHungerGamesBot` on Telegram.
2. Retrieve API key from https://dinner.sea.com/accounts/token, retrieve the token and tell the bot.
3. Play with the available commands:

| Command     | Description |
| ----------- | ----------- |
| `/start`    | Begin chatting with the bot. If it is your first time, the bot will prompt for your key.|
| `/key`      | Check if your key is remembered by the bot.|
| `/newkey`   | Update a new key.        |
| `/menu`     | Check today's menu and place your order.        |
| `/choice`   | Check the current food that you tasked the bot to order.        |
| `/status`   | Check the order status of the entire week.        |
| `/mute`     | Mute or unmute daily reminders.       |
| `/help`     | Introduction and help.        |

4. The bot will not entertain anymore requests 1 minute before `12.30pm`, and proceed to order for you.
5. It will tell you if the order is successful.
6. Remember to collect and eat it. Yumm.

```mermaid
   sequenceDiagram
   Consumer-->API: Book something
   API-->BookingService: Start booking process
   end
   API-->BillingService: Start billing process
```