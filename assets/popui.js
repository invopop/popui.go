
// Global access to the Console UI SDK URL
const CONSOLE_SDK_URL = 'https://cdn.jsdelivr.net/npm/@invopop/console-ui-sdk@0.0.9/index.js';

(function() {
  'use strict';

  // Constants
  const QUERY_SELECTORS = {
    hamburgerButton: '.popui-admin-page-title__wrapper > button',
    sidebar: '.popui-admin-sidebar',
    page: '.popui-admin-page',
    buttonCopy: '.popui-button-copy',
    buttonCopyValue: '[data-copy-value]',
    buttonCopyText: '.popui-button-copy__text',
    buttonCopyPopover: '.popui-button-copy__popover'
  }
  const ACTIVE_MENU_CLASS = 'menu--active'
  const LOADING_CLASS = 'popui-button--loading'
  const POPOVER_VISIBLE_CLASS = 'popui-button-copy__popover--visible'

  // Internal helper functions
  function prepareAccentColor() {
    const urlParams = new URLSearchParams(window.location.search)
    const accentColor = urlParams.get('accent')

    if (accentColor) {
      const root = document.querySelector(':root')
      root.style.setProperty('--workspace-accent-color', accentColor)
    }
  }

  function updateButtonCopyText(input) {
    const container = input.closest(QUERY_SELECTORS.buttonCopy)
    if (!container) return

    const textButton = container.querySelector(QUERY_SELECTORS.buttonCopyText)
    if (!textButton) return

    const value = input.value || input.getAttribute('value') || ''
    const prefixLength = parseInt(input.dataset.prefixLength) || 0
    const suffixLength = parseInt(input.dataset.suffixLength) || 0

    textButton.textContent = formatButtonCopyText(value, prefixLength, suffixLength)
  }

  function formatButtonCopyText(text, prefixLength, suffixLength) {
    if (!text) return ''

    if (!prefixLength && !suffixLength) return text

    if (text.length <= prefixLength + suffixLength) return text

    let result = ''

    if (prefixLength > 0) {
      result += text.substring(0, prefixLength)
    }

    result += '...'

    if (suffixLength > 0) {
      result += text.substring(text.length - suffixLength)
    }

    return result
  }

  // Initialize popui namespace
  const popui = window.popui || {};

  // Public API: Show loading spinner on button
  popui.showButtonSpinner = function(button) {
    const form = button.form || button.closest('form')
    if (form && form.checkValidity()) {
      button.classList.add(LOADING_CLASS)
    }
  };

  // Public API: Copy button value to clipboard
  popui.copyButtonValue = function(button) {
    const container = button.closest(QUERY_SELECTORS.buttonCopy)
    if (!container) return

    const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue)
    if (!input) return

    const value = input.value || input.getAttribute('value') || ''
    if (!value) return

    navigator.clipboard
      .writeText(value)
      .then(() => {
        // Show popover if it exists
        const popover = container.querySelector(QUERY_SELECTORS.buttonCopyPopover)
        if (popover) {
          popover.classList.add(POPOVER_VISIBLE_CLASS)
          setTimeout(() => {
            popover.classList.remove(POPOVER_VISIBLE_CLASS)
          }, 2000)
        }
      })
      .catch((err) => {
        console.error('Failed to copy text: ', err)
      })
  };

  // Public API: Authentication token management
  popui.setAuthToken = function(token) {
    sessionStorage.setItem('_popui_auth_token', token);
  };

  popui.getAuthToken = function() {
    return sessionStorage.getItem('_popui_auth_token');
  };

  popui.clearAuthToken = function() {
    sessionStorage.removeItem('_popui_auth_token');
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
  popui.initAuth = function() {
    if (authInitialized) {
      console.warn('popui.initAuth has already been called');
      return;
    }
    authInitialized = true;

    // HTMX config request handling to add authentication token to same-origin requests
    document.addEventListener('htmx:configRequest', (e) => {
      if (!isSameOrigin(e.detail.path)) return;
      const token = popui.getAuthToken();
      if (token) e.detail.headers['Authorization'] = 'Bearer ' + token;
    });

    // Axios interceptor to add authentication token to same-origin requests
    if (typeof axios !== 'undefined') {
      axios.interceptors.request.use(function (config) {
        if (!isSameOrigin(config.url)) return config;
        const token = popui.getAuthToken();
        if (token) {
          config.headers['Authorization'] = 'Bearer ' + token;
        }
        return config;
      }, function (error) {
        return Promise.reject(error);
      });
    }
  };

  // Assign to window
  window.popui = popui;

  // Prepare accent color from URL parameter, if provided.
  window.onload = function() {
    prepareAccentColor();
  }

  // DOM initialization
  document.addEventListener('DOMContentLoaded', () => {
    // Sidebar
    const button = document.querySelector(QUERY_SELECTORS.hamburgerButton)
    const sidebar = document.querySelector(QUERY_SELECTORS.sidebar)
    const page = document.querySelector(QUERY_SELECTORS.page)

    const showSidebar = (e) => {
      e.stopPropagation()
      sidebar.classList.add(ACTIVE_MENU_CLASS)
      page.classList.add(ACTIVE_MENU_CLASS)
    }

    const hideSidebar = () => {
      sidebar.classList.remove(ACTIVE_MENU_CLASS)
      page.classList.remove(ACTIVE_MENU_CLASS)
    }

    if (button) {
      button.addEventListener('click', showSidebar)
    }

    if (page) {
      page.addEventListener('click', hideSidebar)
    }

    // ButtonCopy
    const containers = document.querySelectorAll(QUERY_SELECTORS.buttonCopy)

    containers.forEach((container) => {
      const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue)
      if (!input) return

      updateButtonCopyText(input)

      input.addEventListener('input', () => {
        updateButtonCopyText(input)
      })
    })
  })

  // Remove any loading class from buttons after browser buttons navigation
  window.addEventListener('visibilitychange', function () {
    if (document.visibilityState !== 'visible') return
    const loadingButtons = document.querySelectorAll(`.${LOADING_CLASS}`)
    loadingButtons.forEach((button) => {
      button.classList.remove(LOADING_CLASS)
    })
  })
})();
