import './style.css';
import './app.css';

import {Calculator} from '../wailsjs/go/main/App';

const textInput = document.getElementById("string-input");
const calculateButton = document.getElementById("calculate-button");
const responseText = document.getElementById("response-text");

async function calculate() {
    const expression = textInput.value.trim();

    if (!expression) {
        responseText.textContent = "Enter a calculation first.";
        return;
    }

    try {
        responseText.textContent = String(await Calculator(expression));
    } catch (error) {
        responseText.textContent = error instanceof Error ? error.message : String(error);
    }
}

calculateButton.addEventListener("click", calculate);
textInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        calculate();
    }
});
