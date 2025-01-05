class Toast {
  /**
   * A class representing a Toast notification.
   * @param level {("info"|"success"|"warning"|"danger")}
   * @param message { string }
   */
  constructor(level, message) {
    this.level = level;
    this.message = message;
  }

  /**
   * Makes the toast container element. A button containing the entire notification.
   * @returns {HTMLButtonElement}
   */
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

  /**
   * Makes the element containing the body of the toast notification.
   * @returns {HTMLSpanElement}
   */
  #makeToastContentElement() {
    const messageContainer = document.createElement("span");
    messageContainer.textContent = this.message;
    return messageContainer;
  }

  /**
   * Presents the toast notification at the end of the given container.
   * @param containerQuerySelector {string} a CSS query selector identifying the container for all toasts.
   */
  show(containerQuerySelector = "#toast-container") {
    const toast = this.#makeToastContainerButton();
    const toastContent = this.#makeToastContentElement();
    toast.appendChild(toastContent);

    const toastContainer = document.querySelector(containerQuerySelector);
    toastContainer.appendChild(toast);
  }
}

document.body.addEventListener("makeToast", onMakeToast);

/**
 * Presents a toast notification when the `makeToast` event is triggered
 * @param e {{detail: {level: string, message: string}}}
 */
function onMakeToast(e) {
  console.log(e);
  const toast = new Toast(e.detail.level, e.detail.message);
  toast.show();
}
