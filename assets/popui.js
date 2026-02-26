// Global access to the Console UI SDK URL
const CONSOLE_SDK_URL =
  "https://cdn.jsdelivr.net/npm/@invopop/console-ui-sdk@0.0.10/index.js";

(function () {
  "use strict";

  // Constants
  const QUERY_SELECTORS = {
    hamburgerButton: ".popui-admin-page-title__wrapper > button",
    sidebar: ".popui-admin-sidebar",
    page: ".popui-admin-page",
    buttonCopy: "[data-button-copy]",
    buttonCopyValue: "[data-copy-value]",
    buttonCopyText: "[data-copy-text]",
    buttonCopyPopover: "[data-copy-popover]",
  };
  const ACTIVE_MENU_CLASS = "menu--active";
  const LOADING_CLASS = "popui-button--loading";
  const POPOVER_VISIBLE_CLASS = "popui-button-copy__popover--visible";

  // Internal helper functions
  function prepareAccentColor() {
    const urlParams = new URLSearchParams(window.location.search);
    let accentColor = urlParams.get("accent");

    if (!accentColor) {
      const el = document.querySelector("[data-accent-color]");
      if (el) {
        accentColor = el.dataset.accentColor;
      }
    }

    if (accentColor) {
      const root = document.querySelector(":root");
      root.style.setProperty("--workspace-accent-color", accentColor);
      root.style.setProperty("--color-base-accent", accentColor);
    }
  }

  function updateButtonCopyText(input) {
    const container = input.closest(QUERY_SELECTORS.buttonCopy);
    if (!container) return;

    const textButton = container.querySelector(QUERY_SELECTORS.buttonCopyText);
    if (!textButton) return;

    const value = input.value || input.getAttribute("value") || "";
    const prefixLength = parseInt(input.dataset.prefixLength) || 0;
    const suffixLength = parseInt(input.dataset.suffixLength) || 0;

    textButton.textContent = formatButtonCopyText(
      value,
      prefixLength,
      suffixLength,
    );
  }

  function formatButtonCopyText(text, prefixLength, suffixLength) {
    if (!text) return "";

    if (!prefixLength && !suffixLength) return text;

    if (text.length <= prefixLength + suffixLength) return text;

    let result = "";

    if (prefixLength > 0) {
      result += text.substring(0, prefixLength);
    }

    result += "...";

    if (suffixLength > 0) {
      result += text.substring(text.length - suffixLength);
    }

    return result;
  }

  // Initialize popui namespace
  const popui = window.popui || {};

  // Public API: Show loading spinner on button
  popui.showButtonSpinner = function (button) {
    const form = button.form || button.closest("form");
    if (form && form.checkValidity()) {
      button.classList.add(LOADING_CLASS);
    }
  };

  // Public API: Copy button value to clipboard
  popui.copyButtonValue = function (button) {
    const container = button.closest(QUERY_SELECTORS.buttonCopy);
    if (!container) return;

    const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue);
    if (!input) return;

    const value = input.value || input.getAttribute("value") || "";
    if (!value) return;

    navigator.clipboard
      .writeText(value)
      .then(() => {
        // Show popover if it exists
        const popover = container.querySelector(
          QUERY_SELECTORS.buttonCopyPopover,
        );
        if (popover) {
          popover.classList.add(POPOVER_VISIBLE_CLASS);
          setTimeout(() => {
            popover.classList.remove(POPOVER_VISIBLE_CLASS);
          }, 2000);
        }
      })
      .catch((err) => {
        console.error("Failed to copy text: ", err);
      });
  };

  // Public API: Copy code block content to clipboard
  popui.copyCodeBlock = function (button) {
    const pre = button.closest("pre");
    if (!pre) return;

    const code = pre.querySelector("code");
    if (!code) return;

    navigator.clipboard
      .writeText(code.textContent)
      .then(() => {
        const duplicateIcon = button.querySelector(
          "[data-copy-icon-duplicate]",
        );
        const successIcon = button.querySelector("[data-copy-icon-success]");
        if (duplicateIcon && successIcon) {
          duplicateIcon.classList.add("hidden");
          successIcon.classList.remove("hidden");
          setTimeout(() => {
            duplicateIcon.classList.remove("hidden");
            successIcon.classList.add("hidden");
          }, 2000);
        }
      })
      .catch((err) => {
        console.error("Failed to copy code: ", err);
      });
  };

  // Public API: Authentication token management
  popui.setAuthToken = function (token) {
    sessionStorage.setItem("_popui_auth_token", token);
  };

  popui.getAuthToken = function () {
    return sessionStorage.getItem("_popui_auth_token");
  };

  popui.clearAuthToken = function () {
    sessionStorage.removeItem("_popui_auth_token");
  };

  // Helper function to check if URL is same origin
  function isSameOrigin(url) {
    if (!url) return true; // Relative URLs are same origin
    try {
      const requestUrl = new URL(url, window.location.origin);
      return requestUrl.origin === window.location.origin;
    } catch (e) {
      return true; // If parsing fails, assume relative URL
    }
  }

  // Track if auth has been initialized to prevent duplicate listeners
  let authInitialized = false;

  // Initialize authentication interceptors
  popui.initAuth = function () {
    if (authInitialized) {
      console.warn("popui.initAuth has already been called");
      return;
    }
    authInitialized = true;

    // HTMX config request handling to add authentication token to same-origin requests
    document.addEventListener("htmx:configRequest", (e) => {
      if (!isSameOrigin(e.detail.path)) return;
      const token = popui.getAuthToken();
      if (token) e.detail.headers["Authorization"] = "Bearer " + token;
    });

    // Axios interceptor to add authentication token to same-origin requests
    if (typeof axios !== "undefined") {
      axios.interceptors.request.use(
        function (config) {
          if (!isSameOrigin(config.url)) return config;
          const token = popui.getAuthToken();
          if (token) {
            config.headers["Authorization"] = "Bearer " + token;
          }
          return config;
        },
        function (error) {
          return Promise.reject(error);
        },
      );
    }
  };

  // Inline position styles for toast containers (using inline styles for reliability)
  const toastPositionStyles = {
    "top-left": "top: 1rem; left: 1rem;",
    "top-center": "top: 1rem; left: 50%; transform: translateX(-50%);",
    "top-right": "top: 1rem; right: 1rem;",
    "bottom-left": "bottom: 1rem; left: 1rem;",
    "bottom-center": "bottom: 1rem; left: 50%; transform: translateX(-50%);",
    "bottom-right": "bottom: 1rem; right: 1rem;",
  };

  // Get position from containerId (e.g., 'toast-top-left' -> 'top-left')
  function getPositionFromContainerId(containerId) {
    for (const pos of Object.keys(toastPositionStyles)) {
      if (containerId.includes(pos)) {
        return pos;
      }
    }
    return "bottom-right"; // default
  }

  // Public API: Show toast notification (works with or without Alpine.js)
  popui.showToast = function (options) {
    console.log("popui.showToast called with:", options);
    const opts =
      typeof options === "string"
        ? { variant: options, title: arguments[1], description: arguments[2] }
        : options;

    const {
      variant = "default",
      title = "",
      description = "",
      duration = 3000,
      containerId = "toast-container",
      position,
    } = opts;

    // Vanilla JS implementation (always use this for reliability)
    let container = document.getElementById(containerId);
    if (!container) {
      container = document.createElement("div");
      container.id = containerId;
      const pos = position || getPositionFromContainerId(containerId);
      const posStyle =
        toastPositionStyles[pos] || toastPositionStyles["bottom-right"];
      container.style.cssText = `position: fixed; z-index: 9999; display: flex; flex-direction: column; gap: 0.5rem; pointer-events: none; ${posStyle}`;
      document.body.appendChild(container);
    }

    // Icons from @icons package
    const icons = {
      success:
        '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 16"><path fill="#018551" d="M8 1a7 7 0 1 1 0 14A7 7 0 0 1 8 1"/><path stroke="#fff" stroke-linecap="round" stroke-linejoin="round" d="m5 8.5 2 2 4-5"/></svg>',
      error:
        '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 16"><path fill="#b30b00" d="M8 1a7 7 0 1 1 0 14A7 7 0 0 1 8 1"/><path stroke="#fff" stroke-linecap="round" stroke-linejoin="round" d="m11 5-6 6M5 5l6 6"/></svg>',
      warning:
        '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 16"><path fill="#ec6d1b" d="M6.698 2.28a1.5 1.5 0 0 1 2.604 0l5.416 9.476A1.5 1.5 0 0 1 13.415 14H2.585a1.5 1.5 0 0 1-1.303-2.244z"/><path fill="#fff" stroke="#fff" stroke-width="1.195" d="M8 11.598c.084 0 .152.068.152.152a.153.153 0 0 1-.152.152.153.153 0 0 1-.152-.152c0-.084.069-.152.152-.152Z"/><path stroke="#fff" d="M8 5v5"/></svg>',
      info: '<svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M8 1C11.8659 1 15 4.13406 15 8C15 11.8659 11.8659 15 8 15C4.13406 15 1 11.8659 1 8C1 4.13406 4.13406 1 8 1ZM8 2C4.68635 2 2 4.68635 2 8C2 11.3137 4.68635 14 8 14C11.3137 14 14 11.3137 14 8C14 4.68635 11.3137 2 8 2ZM8.5 12H7.5V7H8.5V12ZM8 4C8.414 4 8.75 4.33675 8.75 4.75C8.75 5.16325 8.414 5.5 8 5.5C7.586 5.5 7.25 5.16325 7.25 4.75C7.25 4.33675 7.586 4 8 4Z" fill="currentColor"/></svg>',
    };

    const toast = document.createElement("div");
    // rounded-lg (8px), gap-2 (8px), py-2.5 (10px), px-3 (12px)
    // Use items-center when no description, items-start when there is description
    const alignClass = description ? "items-start" : "items-center";
    toast.className = `pointer-events-auto flex ${alignClass} gap-2 rounded-lg shadow-lg max-w-sm w-80 bg-background-default-negative py-2.5 px-3 transition-all duration-300`;
    toast.style.cssText = "opacity: 0; transform: translateY(10px);";

    const icon = icons[variant] || "";
    // Info icon uses text-icon-inverse color class
    const iconColorClass = variant === "info" ? "text-icon-inverse" : "";
    toast.innerHTML = `
      ${icon ? `<div class="flex-shrink-0 ${iconColorClass}" style="width: 16px; height: 16px;">${icon}</div>` : ""}
      <div class="flex-1 min-w-0">
        ${title ? `<p class="text-base font-medium text-foreground-inverse">${title}</p>` : ""}
        ${description ? `<p class="text-base text-foreground-inverse-secondary">${description}</p>` : ""}
      </div>
    `;

    const dismiss = () => {
      toast.style.opacity = "0";
      toast.style.transform = "translateY(10px)";
      setTimeout(() => toast.remove(), 300);
    };

    container.appendChild(toast);
    requestAnimationFrame(() => {
      toast.style.opacity = "1";
      toast.style.transform = "translateY(0)";
    });

    if (duration > 0) {
      setTimeout(dismiss, duration);
    }

    return { dismiss };
  };

  // Initialize Alpine.js toast store
  let toastStoreInitialized = false;
  popui.initToastStore = function () {
    if (toastStoreInitialized) return;
    if (typeof Alpine === "undefined") return;

    toastStoreInitialized = true;
    Alpine.store("toasts", {
      items: {},

      add(options) {
        const id = Date.now() + Math.random().toString(36).substr(2, 9);
        const containerId = options.containerId || "toast-container";
        const duration = options.duration ?? 3000;

        if (!this.items[containerId]) {
          this.items[containerId] = [];
        }

        const toast = { id, ...options, show: true };
        this.items[containerId].push(toast);

        if (duration > 0) {
          setTimeout(() => this.remove(containerId, id), duration);
        }

        return {
          dismiss: () => this.remove(containerId, id),
        };
      },

      remove(containerId, id) {
        if (!this.items[containerId]) return;
        const toast = this.items[containerId].find((t) => t.id === id);
        if (toast) {
          toast.show = false;
          setTimeout(() => {
            this.items[containerId] = this.items[containerId].filter(
              (t) => t.id !== id,
            );
          }, 300);
        }
      },

      getItems(containerId) {
        return this.items[containerId] || [];
      },
    });
  };

  // Assign to window
  window.popui = popui;

  // Register Alpine store during alpine:init (before Alpine starts)
  document.addEventListener("alpine:init", () => {
    popui.initToastStore();
  });

  // Prepare accent color from URL parameter, if provided.
  window.onload = function () {
    prepareAccentColor();
  };

  // Toast initialization using MutationObserver (like templui)
  const toastTimers = new Map();

  function initToast(toastEl) {
    const duration = parseInt(toastEl.dataset.popuiToastDuration || "3000");
    const state = {
      timer: null,
      startTime: Date.now(),
      remaining: duration,
      paused: false,
    };
    toastTimers.set(toastEl, state);

    // Animate in
    requestAnimationFrame(() => {
      toastEl.style.opacity = "1";
      toastEl.style.transform = "translateY(0)";
    });

    // Auto-dismiss after duration
    if (duration > 0) {
      state.timer = setTimeout(() => dismissToast(toastEl), duration);
    }

    // Pause on hover
    toastEl.addEventListener("mouseenter", () => {
      const s = toastTimers.get(toastEl);
      if (!s || s.paused) return;
      clearTimeout(s.timer);
      s.remaining = s.remaining - (Date.now() - s.startTime);
      s.paused = true;
    });

    // Resume on mouse leave
    toastEl.addEventListener("mouseleave", () => {
      const s = toastTimers.get(toastEl);
      if (!s || !s.paused || s.remaining <= 0) return;
      s.startTime = Date.now();
      s.paused = false;
      s.timer = setTimeout(() => dismissToast(toastEl), s.remaining);
    });
  }

  function dismissToast(toastEl) {
    toastTimers.delete(toastEl);
    toastEl.style.opacity = "0";
    toastEl.style.transform = "translateY(1rem)";
    setTimeout(() => toastEl.remove(), 300);
  }

  // Handle dismiss button clicks
  document.addEventListener("click", (e) => {
    const dismissBtn = e.target.closest("[data-popui-toast-dismiss]");
    if (dismissBtn) {
      const toast = dismissBtn.closest("[data-popui-toast]");
      if (toast) dismissToast(toast);
    }
  });

  // Watch for new toasts added to the DOM
  const toastObserver = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      mutation.addedNodes.forEach((node) => {
        if (node.nodeType === 1) {
          // Check if the node itself is a toast
          if (node.matches?.("[data-popui-toast]")) {
            console.log("Toast detected (direct):", node);
            initToast(node);
          }
          // Also check for toasts inside the added node (for HTMX fragments)
          const nestedToasts = node.querySelectorAll?.("[data-popui-toast]");
          if (nestedToasts) {
            nestedToasts.forEach((toast) => {
              console.log("Toast detected (nested):", toast);
              initToast(toast);
            });
          }
        }
      });
    });
  });

  // Start observing once DOM is ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", () => {
      toastObserver.observe(document.body, { childList: true, subtree: true });
    });
  } else {
    toastObserver.observe(document.body, { childList: true, subtree: true });
  }

  // DOM initialization
  document.addEventListener("DOMContentLoaded", () => {
    // Sidebar
    const button = document.querySelector(QUERY_SELECTORS.hamburgerButton);
    const sidebar = document.querySelector(QUERY_SELECTORS.sidebar);
    const page = document.querySelector(QUERY_SELECTORS.page);

    const showSidebar = (e) => {
      e.stopPropagation();
      sidebar.classList.add(ACTIVE_MENU_CLASS);
      page.classList.add(ACTIVE_MENU_CLASS);
    };

    const hideSidebar = () => {
      sidebar.classList.remove(ACTIVE_MENU_CLASS);
      page.classList.remove(ACTIVE_MENU_CLASS);
    };

    if (button) {
      button.addEventListener("click", showSidebar);
    }

    if (page) {
      page.addEventListener("click", hideSidebar);
    }

    // ButtonCopy
    const containers = document.querySelectorAll(QUERY_SELECTORS.buttonCopy);

    containers.forEach((container) => {
      const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue);
      if (!input) return;

      updateButtonCopyText(input);

      input.addEventListener("input", () => {
        updateButtonCopyText(input);
      });
    });
  });

  // Remove any loading class from buttons after browser buttons navigation
  window.addEventListener("visibilitychange", function () {
    if (document.visibilityState !== "visible") return;
    const loadingButtons = document.querySelectorAll(`.${LOADING_CLASS}`);
    loadingButtons.forEach((button) => {
      button.classList.remove(LOADING_CLASS);
    });
  });
})();
