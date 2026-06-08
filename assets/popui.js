
// Global access to the Console UI SDK URL
const CONSOLE_SDK_URL = 'https://cdn.jsdelivr.net/npm/@invopop/console-ui-sdk@0.0.10/index.js';

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
    let accentColor = urlParams.get('accent')

    if (!accentColor) {
      const el = document.querySelector('[data-accent-color]')
      if (el) {
        accentColor = el.dataset.accentColor
      }
    }

    if (accentColor) {
      const root = document.querySelector(':root')
      root.style.setProperty('--workspace-accent-color', accentColor)
      root.style.setProperty('--color-base-accent', accentColor)
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
    // Legacy sidebar (deprecated Page component)
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

    // Modern sidebar toggle (App + Sidebar components)
    const popuiSidebar = document.getElementById('popui-sidebar')
    if (popuiSidebar) {
      const openSidebar = () => popuiSidebar.classList.add('popui-sidebar-open')
      const closeSidebar = () => popuiSidebar.classList.remove('popui-sidebar-open')
      const toggleSidebar = () => popuiSidebar.classList.toggle('popui-sidebar-open')
      const mql = window.matchMedia('(min-width: 768px)')

      // Open by default at md+ breakpoint
      if (mql.matches) openSidebar()

      // Mark as ready so the CSS closed-state rule can take effect
      requestAnimationFrame(() => popuiSidebar.classList.add('popui-sidebar-ready'))

      // Toggle button (open sidebar)
      document.querySelectorAll('[data-sidebar-toggle]').forEach((btn) => {
        btn.addEventListener('click', toggleSidebar)
      })

      // Hide button (close sidebar)
      document.querySelectorAll('[data-sidebar-hide]').forEach((btn) => {
        btn.addEventListener('click', closeSidebar)
      })

      // Close on escape (mobile only)
      document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !mql.matches) closeSidebar()
      })

      // Close when clicking outside the sidebar (mobile only)
      document.addEventListener('click', (e) => {
        if (mql.matches) return
        if (!popuiSidebar.classList.contains('popui-sidebar-open')) return
        if (!popuiSidebar.contains(e.target) && !e.target.closest('[data-sidebar-toggle]')) {
          closeSidebar()
        }
      })

      // Auto-open/close when crossing the md breakpoint
      mql.addEventListener('change', (e) => {
        if (e.matches) openSidebar()
        else closeSidebar()
      })
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

  // Polyfill for anchor positioning on browsers that don't support it
  if (!CSS.supports('anchor-name', '--test')) {
    const positionContextMenu = (contextMenu, trigger) => {
      const triggerRect = trigger.getBoundingClientRect()
      const isRightAlign = contextMenu.classList.contains('context-menu-right-align')

      contextMenu.style.position = 'fixed'
      contextMenu.style.top = `${triggerRect.bottom + 8}px`

      if (isRightAlign) {
        contextMenu.style.left = 'auto'
        contextMenu.style.right = `${window.innerWidth - triggerRect.right}px`
      } else {
        contextMenu.style.left = `${triggerRect.left}px`
        contextMenu.style.right = 'auto'
      }
    }

    document.addEventListener('toggle', (e) => {
      const contextMenu = e.target
      if (!contextMenu.matches('[popover].context-menu')) return

      // Find the trigger button
      // First try the standard attribute (works after Alpine binds it)
      let trigger = document.querySelector(`[popovertarget="${contextMenu.id}"]`)

      // If not found, traverse up from the popover to find the button in the same container
      if (!trigger && contextMenu.parentElement) {
        trigger = contextMenu.parentElement.querySelector('button')
      }

      if (!trigger) return

      if (e.newState === 'open') {
        // Position initially
        positionContextMenu(contextMenu, trigger)

        // Update position on scroll
        const updatePosition = () => positionContextMenu(contextMenu, trigger)
        window.addEventListener('scroll', updatePosition, true)
        window.addEventListener('resize', updatePosition)

        // Clean up listeners when context menu closes
        contextMenu.addEventListener('toggle', function cleanup(e) {
          if (e.newState === 'closed') {
            window.removeEventListener('scroll', updatePosition, true)
            window.removeEventListener('resize', updatePosition)
            contextMenu.removeEventListener('toggle', cleanup)
          }
        })
      }
    }, true)
  }
})();

// Alpine controller for the DropdownSelect component. Single mode submits on
// click and closes the popover. Multiple mode toggles the value in/out of
// the selection set without closing; the form is submitted once on
// popover close, only if the selection actually changed. Auto-open pops
// the dropdown after mount — used by callers that spawn the chip via a
// parent menu and want the options visible without a second click.
document.addEventListener('alpine:init', () => {
  if (!window.Alpine) return
  Alpine.data('dropdownSelect', (init) => ({
    values: (init && init.values) || [],
    multiple: !!(init && init.multiple),
    multipleLabel: (init && init.multipleLabel) || 'items',
    name: (init && init.name) || '',
    initial: '',
    init() {
      this.initial = JSON.stringify(this.values)
      if (init && init.autoOpen) {
        this.$nextTick(() => {
          const trigger = this.$refs.trigger
          if (trigger && typeof trigger.click === 'function') trigger.click()
        })
      }
    },
    submitForm() {
      const form = this.$root.closest('form')
      if (form && typeof form.requestSubmit === 'function') form.requestSubmit()
    },
    toggle(v) {
      if (!this.multiple) {
        this.values = [v]
        const pop = this.$root.querySelector('[popover]')
        if (pop && typeof pop.hidePopover === 'function') pop.hidePopover()
        this.initial = JSON.stringify(this.values)
        // Defer to $nextTick so Alpine's reactive `x-for` has rendered the
        // new hidden inputs before HTMX serialises the form. Without this
        // the form submits with stale (empty) input set and the URL ends
        // up without the new value.
        this.$nextTick(() => this.submitForm())
        return
      }
      if (this.values.includes(v)) {
        this.values = this.values.filter((x) => x !== v)
      } else {
        this.values = [...this.values, v]
      }
      // Multi-mode: submit on every flip so the table refreshes
      // immediately. Sync `initial` first so the onToggle close-handler
      // sees no diff and skips its redundant submit. $nextTick is required
      // so Alpine's reactive `x-for` has rendered the hidden inputs
      // before HTMX reads the form.
      //
      // The consumer should scope the HTMX swap to a region that does NOT
      // include this dropdown (e.g. via hx-target/hx-select on a table
      // wrapper, with hx-swap-oob for any sibling regions that need
      // refreshing too). That way the popover survives the swap and the
      // user can keep ticking without the dropdown re-rendering.
      this.initial = JSON.stringify(this.values)
      this.$nextTick(() => this.submitForm())
    },
    onToggle(e) {
      if (!this.multiple) return
      if (e.newState !== 'closed') return
      const now = JSON.stringify(this.values)
      if (now === this.initial) return
      this.initial = now
      this.submitForm()
    },
  }))
})

// Alpine controller for popui.FilterRow. Tracks the single active filter
// (the row uses a single-filter-at-a-time UX), handles add/remove,
// clears stale values when the active chip changes, and auto-opens the
// chip's value editor on add — matches @invopop/popui Svelte
// FilterChip's justAdded onMount toggle.
//
// Each chip wraps itself in a [data-filter-name=…] div; the controller
// uses that to find the chip's editor and either focus its text input or
// click its popovertarget trigger.
document.addEventListener('alpine:init', () => {
  if (!window.Alpine) return
  Alpine.data('filterRow', (initialActive) => ({
    active: Array.isArray(initialActive) ? [...initialActive] : [],

    // Clear all form fields associated with a filter name. Handles plain
    // text inputs, <select>s, and DropdownSelect — the latter stores its
    // selection in inner Alpine state, not directly in form inputs, so
    // we reach into that scope.
    clearFilter(name) {
      const chip = this.$root.querySelector('[data-filter-name="' + name + '"]')
      if (!chip) return
      const input = chip.querySelector('input[type="text"][name="' + name + '"]')
      if (input) input.value = ''
      const select = chip.querySelector('select[name="' + name + '"]')
      if (select) select.value = ''
      // DropdownSelect inner scope (role=combobox).
      const dropdown = chip.querySelector('[role="combobox"]')
      if (dropdown && window.Alpine) {
        const data = Alpine.$data(dropdown)
        if (data && Array.isArray(data.values)) {
          data.values = []
          data.initial = '[]'
        }
      }
    },

    // Open the chip's value editor right after it appears so the user
    // doesn't need an extra click. requestAnimationFrame + setTimeout
    // fallback handles the timing dance between the closing "+ Filter"
    // popover and the opening chip popover (browsers serialise popover
    // transitions, sometimes refusing showPopover on a popover while
    // another is mid-close).
    autoOpenChip(name) {
      const chip = this.$root.querySelector('[data-filter-name="' + name + '"]')
      if (!chip) return
      const tryOpen = () => {
        const popover = chip.querySelector('[popover]')
        if (popover && popover.matches(':popover-open')) return true
        const trigger = chip.querySelector('button[popovertarget]')
        if (trigger && typeof trigger.click === 'function') { trigger.click(); return true }
        if (popover && typeof popover.showPopover === 'function') {
          try { popover.showPopover(); return true } catch (e) {}
        }
        return false
      }
      const focusInput = () => {
        const txt = chip.querySelector('input[type="text"]')
        if (txt) txt.focus()
      }
      requestAnimationFrame(() => {
        if (tryOpen()) return
        setTimeout(() => { if (!tryOpen()) focusInput() }, 0)
      })
    },

    add(name) {
      // Single-filter model: replace whatever chip was active before and
      // null out the leaving chip's values so they don't leak into the
      // URL on the next submit.
      const previous = this.active.filter((n) => n !== name)
      previous.forEach((n) => this.clearFilter(n))
      this.active = [name]
      this.$nextTick(() => this.autoOpenChip(name))
    },

    remove(name) {
      this.active = this.active.filter((n) => n !== name)
      this.clearFilter(name)
      // Defer the submit to $nextTick so Alpine's reactive <template x-for>
      // inside a DropdownSelect chip has time to remove the hidden input(s)
      // for the cleared values before HTMX serialises the form. Without
      // this the removed value still lives in the DOM at submit time and
      // the URL keeps the filter the user just dismissed — mirrors the
      // deferral dropdownSelect.toggle() already uses for the add case.
      this.$nextTick(() => {
        if (this.$root && typeof this.$root.requestSubmit === 'function') {
          this.$root.requestSubmit()
        }
      })
    },

    isActive(name) {
      return this.active.includes(name)
    },

    available(name) {
      return !this.active.includes(name)
    },
  }))
})

