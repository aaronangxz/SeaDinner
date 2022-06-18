package handlers

import "testing"

func TestMakeHelpResponse(t *testing.T) {
	expected := "*Welcome to SeaHungerGamesBot!*\n\n" +
		"The goal of my existence is to help you snatch that dinner in milliseconds. And also we all know that you are too lazy to open up SeaTalk.\n\n" +
		"*Get started*\n" +
		"1. /key to tell me your Sea API key. This is important because without the key, I'm basically useless. When you refresh your key, remember to let me know in /newkey\n" +
		"2. /menu to browse through the dishes, and tap the button below to snatch. There are also options to choose a random dish or skip ordering. Do take note that if you choose to skip, I will remember that and stop ordering forever until you tell me to do so again.\n" +
		"3. /choice to check the current dish I'm tasked to order.\n" +
		"4. /status to see what you have ordered this week, and the order status.\n" +
		"5. /mute to stop receiving morning reminders. Not recommended tho!\n\n" +
		"*Features*\n" +
		"1. I will send you a daily reminder at 10.30am (If you never mute or skip order on that day). Order can be altered easily from the quick options:\n" +
		"🎲 to order a random dish\n" +
		"🙅 to stop ordering\n" +
		"2. At 12.29pm, I will no longer entertain your requests, because I have better things to do! Don't even think about last minute changes.\n" +
		"3. At 12.30pm sharp, I will begin to order your precious food.\n" +
		"4. It is almost guranteed that I can order it in less than 500ms. Will drop you a message too!\n\n" +
		"*Disclaimer*\n" +
		"By using my services, you agree to let me store your API key. However, not to worry! Your key is encrypted with AES-256, it's very unlikely that it will be stolen.\n\n" +
		"*Contribute*\n" +
		"If you see or encounter any bugs, or if there's any feature / improvement that you have in mind, feel free to open an Issue / Pull Request at https://github.com/aaronangxz/SeaDinner\n\n" +
		"Thank you and happy eating!😋"
	tests := []struct {
		name string
		want string
	}{
		{
			"HappyCase",
			expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeHelpResponse(); got != tt.want {
				t.Errorf("MakeHelpResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
