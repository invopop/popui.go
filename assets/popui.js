
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

// Alpine controller for popui.FilterRow's inline option editor — the
// always-open, trigger-less option list that replaces the old DropdownSelect
// for coloured/multi filter fields. Owns the selected `values`, the arrow-key
// highlight (`activeIndex`), and the keyboard handlers (Down/Up move the
// highlight, Enter toggles the highlighted option, Esc blurs). Selection is
// emitted as reactive hidden inputs by the template; every change submits the
// form on $nextTick so the x-for inputs exist before HTMX serialises.
document.addEventListener('alpine:init', () => {
  if (!window.Alpine) return
  Alpine.data('filterOptionList', (init) => ({
    values: (init && init.values) || [],
    multiple: !!(init && init.multiple),
    multipleLabel: (init && init.multipleLabel) || 'items',
    name: (init && init.name) || '',
    optionValues: (init && init.optionValues) || [],
    activeIndex: -1,
    open: false,
    initial: '',
    init() {
      this.initial = JSON.stringify(this.values)
      // Arrow/Enter/Space/Esc are driven from the document level: focusing the
      // list can be undone by the closing "+ Filter" popover restoring focus to
      // its invoker (the "+" button) or <body>. When CLOSED, only the focused
      // combobox control opens (on Down/Enter/Space) — keys aren't hijacked
      // globally. When OPEN, keys act while focus is in the control, lost to
      // <body>, or anywhere else in the same filter row (e.g. the "+" button);
      // keys aimed at controls outside the row and typing in text filters are
      // left alone.
      this._onKeydown = (e) => {
        if (!this.$root || this.$root.offsetParent === null) return // chip hidden
        if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp' && e.key !== 'Enter' && e.key !== 'Escape' && e.key !== ' ') return
        const a = document.activeElement
        const editable = a && (a.tagName === 'INPUT' || a.tagName === 'TEXTAREA' || a.tagName === 'SELECT' || a.isContentEditable)
        if (editable && !this.$root.contains(a)) return // don't hijack typing elsewhere
        if (!this.open) {
          if (this.$root.contains(a) && (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ')) {
            e.preventDefault()
            this.openPanel()
          }
          return
        }
        const form = this.$root.closest('form')
        if (a && a !== document.body && !this.$root.contains(a) && !(form && form.contains(a))) return // focus outside the row
        this.onKeydown(e)
      }
      document.addEventListener('keydown', this._onKeydown, true)
      // A click outside the chip closes the options panel (summary box stays).
      this._onDocClick = (e) => {
        if (!this.open) return
        if (this.$root && this.$root.contains(e.target)) return
        this.closePanel()
      }
      document.addEventListener('click', this._onDocClick, true)
    },
    destroy() {
      if (this._onKeydown) document.removeEventListener('keydown', this._onKeydown, true)
      if (this._onDocClick) document.removeEventListener('click', this._onDocClick, true)
    },
    submitForm() {
      const form = this.$root.closest('form')
      if (form && typeof form.requestSubmit === 'function') form.requestSubmit()
    },
    openPanel() {
      this.open = true
      if (this.activeIndex < 0 && this.optionValues.length) this.activeIndex = 0
    },
    closePanel() {
      this.open = false
      this.activeIndex = -1
    },
    togglePanel() {
      this.open ? this.closePanel() : this.openPanel()
    },
    move(delta) {
      const n = this.optionValues.length
      if (!n) return
      const base = this.activeIndex < 0 ? (delta > 0 ? -1 : 0) : this.activeIndex
      this.activeIndex = ((base + delta) % n + n) % n
    },
    onKeydown(e) {
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        this.move(1)
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        this.move(-1)
      } else if (e.key === 'Enter' || e.key === ' ') {
        // Enter/Space toggle the highlighted option. preventDefault stops a
        // focused option <button> from also firing its native click (double
        // toggle) and stops Space from scrolling the page.
        e.preventDefault()
        if (this.activeIndex >= 0) this.toggle(this.optionValues[this.activeIndex])
      } else if (e.key === 'Escape') {
        this.closePanel()
      }
    },
    // choose moves the highlight to a row and toggles it — the single entry
    // point for a mouse click on a row, so click and Enter stay in sync (the
    // rows are tabindex=-1 and never hold the keyboard cursor themselves).
    choose(i) {
      if (i < 0 || i >= this.optionValues.length) return
      this.activeIndex = i
      this.toggle(this.optionValues[i])
    },
    toggle(v) {
      if (!this.multiple) {
        this.values = [v]
      } else if (this.values.includes(v)) {
        this.values = this.values.filter((x) => x !== v)
      } else {
        this.values = [...this.values, v]
      }
      this.initial = JSON.stringify(this.values)
      // Defer so Alpine's reactive x-for hidden inputs render before HTMX
      // serialises the form (mirrors dropdownSelect.toggle).
      this.$nextTick(() => this.submitForm())
    },
  }))
})

// Alpine controller for popui.FilterRow. Tracks the active filters
// (multi-filter UX: several apply at once), handles add/remove, clears a
// filter's stale values when its chip is removed, and auto-opens a chip's
// value editor on add.
//
// Each chip wraps itself in a [data-filter-name=…] div; the controller
// uses that to find the chip's editor and either focus its text input or
// click its popovertarget trigger.
document.addEventListener('alpine:init', () => {
  if (!window.Alpine) return
  Alpine.data('filterRow', (initialActive, allNames) => ({
    active: Array.isArray(initialActive) ? [...initialActive] : [],
    // Every filterable field name, in declaration order. Used to decide when
    // the "+ Filter" add button is still useful — once every field is active
    // there is nothing left to add, so the button hides.
    all: Array.isArray(allNames) ? [...allNames] : [],

    // Clear all form fields associated with a filter name. Handles plain
    // text inputs, <select>s, and DropdownSelect — the latter stores its
    // selection in inner Alpine state, not directly in form inputs, so
    // we reach into that scope.
    clearFilter(name) {
      // Resolve the form even when called from the nested "+ Filter" menu
      // scope (there Alpine's $root is the menu div, not the form).
      const root = (this.$root.closest && this.$root.closest('form')) || this.$root
      const chip = root.querySelector('[data-filter-name="' + name + '"]')
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
      // Inline option-list editor (data-filter-options) — same reset, plus
      // the arrow-key highlight so a re-added chip starts clean.
      const optionList = chip.querySelector('[data-filter-options]')
      if (optionList && window.Alpine) {
        const listData = Alpine.$data(optionList)
        if (listData && Array.isArray(listData.values)) {
          listData.values = []
          listData.initial = '[]'
          listData.activeIndex = -1
        }
      }
      // Date-range calendar editor — reset the selected range.
      const calendar = chip.querySelector('[data-filter-calendar]')
      if (calendar && window.Alpine) {
        const calData = Alpine.$data(calendar)
        if (calData && typeof calData.clear === 'function') calData.clear()
      }
    },

    // Open the chip's value editor right after it appears so the user
    // doesn't need an extra click. requestAnimationFrame + setTimeout
    // fallback handles the timing dance between the closing "+ Filter"
    // popover and the opening chip popover (browsers serialise popover
    // transitions, sometimes refusing showPopover on a popover while
    // another is mid-close).
    autoOpenChip(name) {
      // Resolve the form even when called from the nested "+ Filter" menu
      // scope (there Alpine's $root is the menu div, not the form), so the
      // chip lookup below actually finds the chip.
      const root = (this.$root.closest && this.$root.closest('form')) || this.$root
      const chip = root.querySelector('[data-filter-name="' + name + '"]')
      if (!chip) return
      // Inline option-list editor: open its panel so the options auto-display
      // on add (its arrow/Enter keys are handled by a document-level listener,
      // so no fragile focus dance is needed).
      const optionList = chip.querySelector('[data-filter-options]')
      if (optionList) {
        if (window.Alpine) {
          const d = Alpine.$data(optionList)
          if (d && typeof d.openPanel === 'function') d.openPanel()
        }
        // Move focus onto the combobox control so it lands on the dropdown —
        // not the "+" button, where the closing "+ Filter" popover would
        // otherwise restore focus (its invoker). The re-assert covers a
        // restoration that lands a beat later; it's a single retry, not a loop
        // (a tight loop fights the popover and never settles).
        if (typeof optionList.focus === 'function') {
          optionList.focus()
          setTimeout(() => {
            if (document.activeElement !== optionList && !optionList.contains(document.activeElement)) {
              optionList.focus()
            }
          }, 50)
        }
        return
      }
      // Date-range calendar editor: open its panel on add, same as the
      // inline option list.
      const calendar = chip.querySelector('[data-filter-calendar]')
      if (calendar) {
        if (window.Alpine) {
          const d = Alpine.$data(calendar)
          if (d && typeof d.openPanel === 'function') d.openPanel()
        }
        if (typeof calendar.focus === 'function') calendar.focus()
        return
      }
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
      // Multi-filter: append the field (don't replace the others) so several
      // filters apply at once and the chips lay out left-to-right in the order
      // they were added. CSS flex `order` on each chip (= its index in
      // `active`) renders them in add-order regardless of input declaration
      // order; the "+" add button is pinned last.
      if (!this.active.includes(name)) this.active = [...this.active, name]
      this.$nextTick(() => this.autoOpenChip(name))
    },

    // orderOf returns a chip's flex order — its position in the add-order
    // `active` list — so chips lay out left-to-right in the order added.
    orderOf(name) {
      const i = this.active.indexOf(name)
      return i < 0 ? 0 : i
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

    hasActive() {
      return this.active.length > 0
    },

    // hasAvailable reports whether any field is still inactive (i.e. the
    // "+ Filter" add menu has something to offer). False once every field is
    // active, which hides the add button.
    hasAvailable() {
      return this.all.some((n) => !this.active.includes(n))
    },
  }))

  // Date-range calendar for the FilterRow's Range chip. Dual-month grid with
  // start/middle/end range selection, prev/next nav, and a preset rail
  // (this/last week, month, quarter + custom). Selection is exposed as
  // `rangeValue` ("YYYY-MM-DD..YYYY-MM-DD") which the chip binds to a hidden
  // input; completing a range (or picking a preset) submits the form.
  Alpine.data('rangeCalendar', (init) => ({
    name: (init && init.name) || '',
    open: false,
    preset: 'custom',
    from: (init && init.from) || null,
    to: (init && init.to) || null,
    viewY: 2000,
    viewM: 0,
    dows: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
    presets: [
      { key: 'thisWeek', label: 'This Week' },
      { key: 'lastWeek', label: 'Last Week' },
      { key: 'thisMonth', label: 'This month' },
      { key: 'lastMonth', label: 'Last month' },
      { key: 'thisQuarter', label: 'This quarter' },
      { key: 'lastQuarter', label: 'Last quarter' },
      { key: 'custom', label: 'Custom' },
    ],
    init() {
      const base = this.from ? this.parse(this.from) : new Date()
      this.viewY = base.getFullYear()
      this.viewM = base.getMonth()
      this._onDocClick = (e) => {
        if (this.open && this.$root && !this.$root.contains(e.target)) this.open = false
      }
      document.addEventListener('click', this._onDocClick, true)
    },
    destroy() {
      document.removeEventListener('click', this._onDocClick, true)
    },
    parse(s) {
      const p = String(s).split('-').map(Number)
      return new Date(p[0], p[1] - 1, p[2])
    },
    iso(dt) {
      const m = String(dt.getMonth() + 1).padStart(2, '0')
      const d = String(dt.getDate()).padStart(2, '0')
      return dt.getFullYear() + '-' + m + '-' + d
    },
    isToday(dt) {
      const t = new Date()
      return dt.getFullYear() === t.getFullYear() && dt.getMonth() === t.getMonth() && dt.getDate() === t.getDate()
    },
    monthGrid(y, m) {
      const first = new Date(y, m, 1)
      const cur = new Date(y, m, 1 - first.getDay())
      const weeks = []
      for (let w = 0; w < 6; w++) {
        const week = []
        for (let d = 0; d < 7; d++) {
          week.push({ day: cur.getDate(), iso: this.iso(cur), outside: cur.getMonth() !== m, today: this.isToday(cur) })
          cur.setDate(cur.getDate() + 1)
        }
        weeks.push(week)
      }
      return weeks
    },
    monthLabel(y, m) {
      return new Date(y, m, 1).toLocaleString('en-US', { month: 'long', year: 'numeric' })
    },
    get monthsView() {
      const r = new Date(this.viewY, this.viewM + 1, 1)
      return [
        { label: this.monthLabel(this.viewY, this.viewM), weeks: this.monthGrid(this.viewY, this.viewM) },
        { label: this.monthLabel(r.getFullYear(), r.getMonth()), weeks: this.monthGrid(r.getFullYear(), r.getMonth()) },
      ]
    },
    prev() {
      const d = new Date(this.viewY, this.viewM - 1, 1)
      this.viewY = d.getFullYear(); this.viewM = d.getMonth()
    },
    next() {
      const d = new Date(this.viewY, this.viewM + 1, 1)
      this.viewY = d.getFullYear(); this.viewM = d.getMonth()
    },
    dayState(iso) {
      if (!this.from) return null
      if (!this.to) return iso === this.from ? 'start' : null
      if (iso === this.from) return 'start'
      if (iso === this.to) return 'end'
      if (iso > this.from && iso < this.to) return 'middle'
      return null
    },
    selectDay(iso) {
      this.preset = 'custom'
      if (!this.from || (this.from && this.to)) {
        this.from = iso; this.to = null
        return
      }
      if (iso < this.from) { this.to = this.from; this.from = iso }
      else { this.to = iso }
      this.apply()
    },
    setPreset(key) {
      this.preset = key
      if (key === 'custom') return
      const t = new Date(); t.setHours(0, 0, 0, 0)
      const sow = (d) => { const x = new Date(d); x.setDate(x.getDate() - x.getDay()); return x }
      let from, to
      if (key === 'thisWeek') { from = sow(t); to = new Date(from); to.setDate(to.getDate() + 6) }
      else if (key === 'lastWeek') { to = sow(t); to.setDate(to.getDate() - 1); from = new Date(to); from.setDate(from.getDate() - 6) }
      else if (key === 'thisMonth') { from = new Date(t.getFullYear(), t.getMonth(), 1); to = new Date(t.getFullYear(), t.getMonth() + 1, 0) }
      else if (key === 'lastMonth') { from = new Date(t.getFullYear(), t.getMonth() - 1, 1); to = new Date(t.getFullYear(), t.getMonth(), 0) }
      else if (key === 'thisQuarter') { const q = Math.floor(t.getMonth() / 3); from = new Date(t.getFullYear(), q * 3, 1); to = new Date(t.getFullYear(), q * 3 + 3, 0) }
      else if (key === 'lastQuarter') { let q = Math.floor(t.getMonth() / 3) - 1, y = t.getFullYear(); if (q < 0) { q = 3; y-- } from = new Date(y, q * 3, 1); to = new Date(y, q * 3 + 3, 0) }
      else return
      this.from = this.iso(from); this.to = this.iso(to)
      this.viewY = from.getFullYear(); this.viewM = from.getMonth()
      this.apply()
    },
    fmt(iso) {
      return this.parse(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
    },
    get summary() {
      if (this.from && this.to) return this.fmt(this.from) + ' – ' + this.fmt(this.to)
      if (this.from) return this.fmt(this.from) + ' – …'
      return ''
    },
    get rangeValue() {
      if (this.from && this.to) return this.from + '..' + this.to
      return ''
    },
    togglePanel() { this.open = !this.open },
    openPanel() { this.open = true },
    clear() { this.from = null; this.to = null; this.preset = 'custom' },
    apply() {
      this.open = false
      const form = this.$root && this.$root.closest && this.$root.closest('form')
      if (form && typeof form.requestSubmit === 'function') {
        this.$nextTick(() => form.requestSubmit())
      }
    },
  }))
})

// Alpine controller for popui.Drawer — a fixed-position floating side
// panel that overlays one edge of the viewport. Non-blocking (no
// backdrop, rest of the app stays interactive), fades in/out on the
// open state, closes on Escape. Stays mounted across HTMX content
// swaps so the caller can re-fill the panel's inner content slot
// without the open/close state flickering.
//
// Open/close protocol: dispatch a `popui-drawer-open` /
// `popui-drawer-close` event on `window` with `detail` set to the
// drawer's ID. Scoped by ID so multiple drawers can coexist on a
// single page without trampling each other.
document.addEventListener('alpine:init', () => {
  if (!window.Alpine) return
  Alpine.data('drawer', (id) => ({
    open: false,
    init() {
      const matches = (e) => e && e.detail === id
      this._onOpen = (e) => { if (matches(e)) this.open = true }
      this._onClose = (e) => { if (matches(e)) this.open = false }
      this._onKeydown = (e) => {
        // Escape only closes the drawer when it's actually open — lets
        // page-level Escape handlers (search bars, popovers) keep
        // working when the drawer is dismissed.
        if (e.key === 'Escape' && this.open) {
          this.open = false
          e.stopPropagation()
        }
      }
      window.addEventListener('popui-drawer-open', this._onOpen)
      window.addEventListener('popui-drawer-close', this._onClose)
      document.addEventListener('keydown', this._onKeydown)

      // Re-broadcast the matching window event whenever `open` flips
      // so external consumers can react to *any* close path (X click,
      // Escape, programmatic), not just the ones that explicitly
      // dispatch the event themselves. Skip equal-value writes to
      // avoid an event echo loop when the listeners above resync the
      // property to the value it already had.
      this.$watch('open', (newVal, oldVal) => {
        if (newVal === oldVal) return
        const name = newVal ? 'popui-drawer-open' : 'popui-drawer-close'
        window.dispatchEvent(new CustomEvent(name, { detail: id }))
      })
    },
    destroy() {
      window.removeEventListener('popui-drawer-open', this._onOpen)
      window.removeEventListener('popui-drawer-close', this._onClose)
      document.removeEventListener('keydown', this._onKeydown)
    },
  }))
})


// Column resizing for popui.Table with Resizable set. Attaches a drag handle
// to the right edge of each non-last header cell of any .popui-table-resizable
// table; dragging sets an inline width/min-width on the th, which the
// table-auto layout respects. Session-local (no persistence). Self-contained
// so consumers don't need their own resize script.
// Leading semicolon guards against ASI: popui.js's prior statement has no
// trailing semicolon, so without this the IIFE would be parsed as a call on
// the previous expression.
;(function () {
  function attachResizers() {
    document.querySelectorAll('.popui-table-resizable thead th').forEach(function (th, idx, all) {
      if (idx === all.length - 1) return; // last column — nothing to drag against
      if (th.querySelector('.popui-table-resizer')) return; // already wired
      var handle = document.createElement('div');
      handle.className = 'popui-table-resizer';
      th.appendChild(handle);
    });
  }
  document.addEventListener('DOMContentLoaded', attachResizers);
  document.addEventListener('htmx:afterSettle', attachResizers);
  document.addEventListener('mousedown', function (e) {
    var handle = e.target.closest('.popui-table-resizer');
    if (!handle) return;
    e.preventDefault();
    var th = handle.closest('th');
    if (!th) return;
    var startX = e.clientX;
    var startWidth = th.offsetWidth;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    function onMove(mv) {
      var next = Math.max(60, startWidth + (mv.clientX - startX));
      th.style.width = next + 'px';
      th.style.minWidth = next + 'px';
    }
    function onUp() {
      document.removeEventListener('mousemove', onMove);
      document.removeEventListener('mouseup', onUp);
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    }
    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
  });
})();
