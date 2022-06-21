package handlers

import "github.com/aaronangxz/SeaDinner/processors"

//MakeHelpResponse Prints out Introduction
func MakeHelpResponse() string {
	txn := processors.App.StartTransaction("make_help_response")
	defer txn.End()
	return "*Welcome to SeaHungerGamesBot!*\n\n" +
		"The goal of my existence is to help you snatch that dinner in milliseconds. And also we all know that you are too lazy to open up SeaTalk.\n\n" +
		"*Important*\n" +
		"I am only good at ordering for the *5SPD* peeps, for now. If you are from other offices, do order from SeaTalk! I am working hard on this :)\n" +
		"Once in a while, SeaTalk will change the order/menu release timings. In such situation, I apologize for not be able to order for you. I will fix myself ASAP to accommodate the changes.\n\n" +
		"*Get started*\n" +
		"1. /key to tell me your Sea API key. This is important because without the key, I'm basically useless. When you refresh your key, remember to let me know in /newkey\n" +
		"2. /menu to browse through the dishes, and tap the button below to snatch. There are also options to choose a random dish or skip ordering. Do take note that if you choose to skip, I will remember that and stop ordering forever until you tell me to do so again.\n" +
		"3. /choice to check the current dish I'm tasked to order.\n" +
		"4. /status to see what you have ordered this week/month, and the order status.\n" +
		"5. /mute to stop receiving morning reminders. Not recommended tho!\n\n" +
		"*Features*\n" +
		"1. I will send you a daily reminder at 10.30am (If you did not mute or skip order on that day). Order can be altered easily from the quick options:\n" +
		"🎲 to order a random dish - excluding vegetarian dishes!\n" +
		"🙅 to stop ordering - I will stop sending reminders to you as well. Choose this option if you are not coming to office.\n" +
		"If you do not change your choice on the subsequent day, I will order the same for you. If you think you will forget, you can plan beforehand by selecting the dish after 12.30pm! Watch out for the **Snatch Tomorrow** buttons\n" +
		"Dish selection will be disabled from Friday afternoon until the next Monday.\n" +
		"2. At 12.29pm, I will no longer entertain your requests, because I have better things to do! Don't even think about last minute changes.\n" +
		"3. At 12.30pm sharp, I will begin to order your precious food.\n" +
		"4. It is almost guaranteed that I can order it in less than 500ms. Will drop you a message too!\n\n" +
		"*Disclaimer*\n" +
		"By using my services, you agree to let me store your API key. However, not to worry! Your key is encrypted with AES-256, it's very unlikely that it will be stolen.\n\n" +
		"*Contribute*\n" +
		"If you see or encounter any bugs, or if there's any feature / improvement that you have in mind, feel free to open an Issue / Pull Request at https://github.com/aaronangxz/SeaDinner\n\n" +
		"Thank you and happy eating!😋"
}
