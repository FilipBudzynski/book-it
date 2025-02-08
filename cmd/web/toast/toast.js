class Toast {
  constructor(level, message) {
    this.level = level;
    this.message = message;
  }

  #makeToastContainerButton() {
    const toastContainer = document.createElement("div");
    toastContainer.classList.add(
      "alert",
      `alert-${this.level}`,
      "transition-opacity",
      "duration-500",
      "ease-in-out",
      "opacity-100",
    );

    setTimeout(() => {
      toastContainer.classList.replace("opacity-100", "opacity-0");
    }, 4500);

    setTimeout(() => {
      toastContainer.remove();
    }, 5000);

    return toastContainer;
  }

  #makeToastContentElement() {
    const messageContainer = document.createElement("span");
    messageContainer.textContent = this.message;
    return messageContainer;
  }

  show(containerQuerySelector = "#toast-container") {
    const toast = this.#makeToastContainerButton();
    const toastContent = this.#makeToastContentElement();
    toast.appendChild(toastContent);

    const toastContainer = document.querySelector(containerQuerySelector);
    toastContainer.appendChild(toast);
  }
}

document.body.addEventListener("makeToast", onMakeToast);

function onMakeToast(e) {
  console.log(e);
  const toast = new Toast(e.detail.level, e.detail.message);
  toast.show();
}
