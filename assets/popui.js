// Global URL for the Console UI SDK.
const CONSOLE_SDK_URL = 'https://cdn.jsdelivr.net/npm/@invopop/console-ui-sdk@0.0.11/index.js';

(function () {
  'use strict';

  const QUERY_SELECTORS = {
    hamburgerButton: '.popui-admin-page-title__wrapper > button',
    sidebar: '.popui-admin-sidebar',
    page: '.popui-admin-page',
    buttonCopy: '[data-button-copy]',
    buttonCopyValue: '[data-copy-value]',
    buttonCopyText: '[data-copy-text]'
  }
  const ACTIVE_MENU_CLASS = 'menu--active'
  const LOADING_CLASS = 'popui-button--loading'

  // ------------------------------------------------------------------
  // Shared helpers
  // ------------------------------------------------------------------

  // Submits the closest enclosing form of an element, when there is one.
  function submitClosestForm(el) {
    const form = el && el.closest && el.closest('form')
    if (form && typeof form.requestSubmit === 'function') form.requestSubmit()
  }

  // Returns the Alpine data scope bound to an element, or null.
  function alpineData(el) {
    return el && window.Alpine ? Alpine.$data(el) : null
  }

  // Returns the selection with a value toggled: single mode replaces it, multiple mode flips it in or out.
  function toggleValue(values, v, multiple) {
    if (!multiple) return [v]
    return values.includes(v) ? values.filter((x) => x !== v) : [...values, v]
  }

  // Shifts a filter chip's dropdown panel left just enough to keep it inside
  // the viewport, retrying on the next frames while the panel has no layout yet.
  function clampPanelX(panel, attempt = 0) {
    if (!panel) return
    const margin = 8
    const cur = parseFloat(panel.style.left) || 0
    const rect = panel.getBoundingClientRect()
    if (!rect.width) {
      if (attempt < 10) requestAnimationFrame(() => clampPanelX(panel, attempt + 1))
      return
    }
    const naturalLeft = rect.left - cur
    const overflow = naturalLeft + rect.width - (window.innerWidth - margin)
    const shift = Math.min(Math.max(overflow, 0), Math.max(naturalLeft - margin, 0))
    const next = shift > 0 ? `${-shift}px` : ''
    if (panel.style.left !== next) panel.style.left = next
  }

  // Applies the workspace accent color from the URL or a data attribute.
  function prepareAccentColor() {
    const urlParams = new URLSearchParams(window.location.search)
    let accentColor = urlParams.get('accent')
    if (!accentColor) {
      const el = document.querySelector('[data-accent-color]')
      if (el) accentColor = el.dataset.accentColor
    }
    if (accentColor) {
      const root = document.querySelector(':root')
      root.style.setProperty('--workspace-accent-color', accentColor)
      root.style.setProperty('--color-base-accent', accentColor)
    }
  }

  // ------------------------------------------------------------------
  // ButtonCopy
  // ------------------------------------------------------------------

  // Renders a ButtonCopy's visible text from its hidden value input.
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

  // Truncates a value to "prefix...suffix" for display.
  function formatButtonCopyText(text, prefixLength, suffixLength) {
    if (!text) return ''
    if (!prefixLength && !suffixLength) return text
    if (text.length <= prefixLength + suffixLength) return text
    let result = ''
    if (prefixLength > 0) result += text.substring(0, prefixLength)
    result += '...'
    if (suffixLength > 0) result += text.substring(text.length - suffixLength)
    return result
  }

  // Populates every ButtonCopy's text and keeps it in sync with its value input.
  function initButtonCopies(root) {
    const scope = root || document
    scope.querySelectorAll(QUERY_SELECTORS.buttonCopy).forEach((container) => {
      const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue)
      if (!input) return
      updateButtonCopyText(input)
      if (input._popuiCopyBound) return
      input._popuiCopyBound = true
      input.addEventListener('input', () => updateButtonCopyText(input))
    })
  }

  // ------------------------------------------------------------------
  // Public API (window.popui)
  // ------------------------------------------------------------------

  const popui = window.popui || {};

  // Shows a loading spinner on a submit button when its form is valid.
  popui.showButtonSpinner = function (button) {
    const form = button.form || button.closest('form')
    if (form && form.checkValidity()) {
      button.classList.add(LOADING_CLASS)
    }
  };

  // Copies a ButtonCopy's value to the clipboard and briefly shows the success icon.
  popui.copyButtonValue = function (button) {
    const container = button.closest(QUERY_SELECTORS.buttonCopy)
    if (!container) return
    const input = container.querySelector(QUERY_SELECTORS.buttonCopyValue)
    if (!input) return
    const value = input.value || input.getAttribute('value') || ''
    if (!value) return
    navigator.clipboard
      .writeText(value)
      .then(() => {
        const dup = container.querySelector('[data-copy-icon-duplicate]')
        const ok = container.querySelector('[data-copy-icon-success]')
        if (!dup || !ok) return
        dup.classList.add('hidden')
        ok.classList.remove('hidden')
        if (container._popuiCopyTimer) clearTimeout(container._popuiCopyTimer)
        container._popuiCopyTimer = setTimeout(() => {
          ok.classList.add('hidden')
          dup.classList.remove('hidden')
        }, 2000)
      })
      .catch((err) => {
        console.error('Failed to copy text: ', err)
      })
  };

  // Toasts: only one is visible at a time and each hides after its data-duration (default 3000ms).
  const TOAST_VISIBLE_CLASS = 'popui-toast--visible'
  const TOAST_DEFAULT_DURATION = 3000
  let activeToast = null
  let activeToastTimer = null

  // Shows a toast by element or id.
  popui.showToast = function (toast) {
    if (typeof toast === 'string') toast = document.getElementById(toast)
    if (!toast) return
    if (activeToastTimer) {
      clearTimeout(activeToastTimer)
      activeToastTimer = null
    }
    if (activeToast && activeToast !== toast) {
      activeToast.classList.remove(TOAST_VISIBLE_CLASS)
    }
    activeToast = toast
    toast.classList.add(TOAST_VISIBLE_CLASS)
    const duration = parseInt(toast.dataset.duration) || TOAST_DEFAULT_DURATION
    activeToastTimer = setTimeout(() => {
      popui.hideToast(toast)
    }, duration)
  };

  // Hides a toast by element or id.
  popui.hideToast = function (toast) {
    if (typeof toast === 'string') toast = document.getElementById(toast)
    if (!toast) return
    toast.classList.remove(TOAST_VISIBLE_CLASS)
    if (activeToast === toast) {
      activeToast = null
      if (activeToastTimer) {
        clearTimeout(activeToastTimer)
        activeToastTimer = null
      }
    }
  };

  // Shows the toast referenced by any clicked element's data-toast-trigger attribute.
  document.addEventListener('click', (e) => {
    const trigger = e.target.closest('[data-toast-trigger]')
    if (!trigger) return
    popui.showToast(trigger.dataset.toastTrigger)
  })

  // Stores the session auth token.
  popui.setAuthToken = function (token) {
    sessionStorage.setItem('_popui_auth_token', token);
  };

  // Returns the session auth token.
  popui.getAuthToken = function () {
    return sessionStorage.getItem('_popui_auth_token');
  };

  // Removes the session auth token.
  popui.clearAuthToken = function () {
    sessionStorage.removeItem('_popui_auth_token');
  };

  // Reports whether a URL shares the page's origin, treating relative or unparsable URLs as same origin.
  function isSameOrigin(url) {
    if (!url) return true;
    try {
      const requestUrl = new URL(url, window.location.origin);
      return requestUrl.origin === window.location.origin;
    } catch (e) {
      return true;
    }
  }

  let authInitialized = false;

  // Installs HTMX and axios interceptors that attach the auth token to same-origin requests.
  popui.initAuth = function () {
    if (authInitialized) {
      console.warn('popui.initAuth has already been called');
      return;
    }
    authInitialized = true;

    document.addEventListener('htmx:configRequest', (e) => {
      if (!isSameOrigin(e.detail.path)) return;
      const token = popui.getAuthToken();
      if (token) e.detail.headers['Authorization'] = 'Bearer ' + token;
    });

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

  window.popui = popui;

  window.onload = function () {
    prepareAccentColor();
  }

  // ------------------------------------------------------------------
  // Input clear buttons
  // ------------------------------------------------------------------

  // Shows or hides a clearable input's clear button based on its value.
  function updateInputClear(input) {
    const btn = input.parentElement ? input.parentElement.querySelector('[data-input-clear]') : null
    if (btn) btn.hidden = input.value === ''
  }

  // Syncs every clearable input's clear button with its current value.
  function initInputClears() {
    document.querySelectorAll('input[data-clearable]').forEach(updateInputClear)
  }

  function onInputClearEvent(e) {
    const t = e.target
    if (t instanceof HTMLInputElement && t.hasAttribute('data-clearable')) updateInputClear(t)
  }

  document.addEventListener('input', onInputClearEvent)
  document.addEventListener('change', onInputClearEvent)

  document.addEventListener('click', (e) => {
    const btn = e.target.closest ? e.target.closest('[data-input-clear]') : null
    if (!btn) return
    const input = btn.parentElement.querySelector('input[data-clearable]')
    if (!input) return
    input.value = ''
    input.dispatchEvent(new Event('input', { bubbles: true }))
    input.dispatchEvent(new Event('change', { bubbles: true }))
    input.focus()
  })

  // Opens the native picker for inputs whose calendar indicator is replaced
  // by the popui calendar button.
  document.addEventListener('click', (e) => {
    const btn = e.target.closest ? e.target.closest('[data-input-picker]') : null
    if (!btn) return
    const input = btn.parentElement.querySelector('input[data-pickable]')
    if (!input) return
    input.focus()
    if (typeof input.showPicker === 'function') {
      try {
        input.showPicker()
      } catch {
        // showPicker can throw (e.g. cross-origin iframe); focus is enough.
      }
    }
  })

  // ------------------------------------------------------------------
  // DOM wiring
  // ------------------------------------------------------------------

  document.addEventListener('DOMContentLoaded', () => {
    // Legacy sidebar used by the deprecated Page component.
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
    if (button) button.addEventListener('click', showSidebar)
    if (page) page.addEventListener('click', hideSidebar)

    // Sidebar toggle for the App + Sidebar components: open at md+ by default, closable on mobile.
    const popuiSidebar = document.getElementById('popui-sidebar')
    if (popuiSidebar) {
      const openSidebar = () => popuiSidebar.classList.add('popui-sidebar-open')
      const closeSidebar = () => popuiSidebar.classList.remove('popui-sidebar-open')
      const toggleSidebar = () => popuiSidebar.classList.toggle('popui-sidebar-open')
      const mql = window.matchMedia('(min-width: 768px)')
      if (mql.matches) openSidebar()
      requestAnimationFrame(() => popuiSidebar.classList.add('popui-sidebar-ready'))
      document.querySelectorAll('[data-sidebar-toggle]').forEach((btn) => {
        btn.addEventListener('click', toggleSidebar)
      })
      document.querySelectorAll('[data-sidebar-hide]').forEach((btn) => {
        btn.addEventListener('click', closeSidebar)
      })
      document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !mql.matches) closeSidebar()
      })
      document.addEventListener('click', (e) => {
        if (mql.matches) return
        if (!popuiSidebar.classList.contains('popui-sidebar-open')) return
        if (!popuiSidebar.contains(e.target) && !e.target.closest('[data-sidebar-toggle]')) {
          closeSidebar()
        }
      })
      mql.addEventListener('change', (e) => {
        if (e.matches) openSidebar()
        else closeSidebar()
      })
    }

    initButtonCopies()
    attachTableResizers()
    initInputClears()
  })

  // Rewires ButtonCopies, table resizers and input clear buttons inserted by HTMX content swaps.
  document.addEventListener('htmx:afterSettle', () => {
    initButtonCopies()
    attachTableResizers()
    initInputClears()
  })

  // Clears button loading spinners when the page becomes visible again.
  window.addEventListener('visibilitychange', function () {
    if (document.visibilityState !== 'visible') return
    document.querySelectorAll(`.${LOADING_CLASS}`).forEach((button) => {
      button.classList.remove(LOADING_CLASS)
    })
  })

  // ------------------------------------------------------------------
  // Anchor-positioning polyfill
  // ------------------------------------------------------------------

  // Positions context-menu popovers under their trigger on browsers without CSS anchor positioning.
  if (!CSS.supports('anchor-name', '--test')) {
    const positionContextMenu = (contextMenu, trigger) => {
      const triggerRect = trigger.getBoundingClientRect()
      const isRightAlign = contextMenu.classList.contains('context-menu-right-align')
      const menuWidth = contextMenu.offsetWidth
      let left = isRightAlign ? triggerRect.right - menuWidth : triggerRect.left
      left = Math.min(left, window.innerWidth - menuWidth - 8)
      left = Math.max(left, 8)
      contextMenu.style.position = 'fixed'
      contextMenu.style.top = `${triggerRect.bottom + 8}px`
      contextMenu.style.left = `${left}px`
      contextMenu.style.right = 'auto'
    }

    document.addEventListener('toggle', (e) => {
      const contextMenu = e.target
      if (!contextMenu.matches('[popover].context-menu')) return
      let trigger = document.querySelector(`[popovertarget="${contextMenu.id}"]`)
      if (!trigger && contextMenu.parentElement) {
        trigger = contextMenu.parentElement.querySelector('button')
      }
      if (!trigger) return
      if (e.newState === 'open') {
        positionContextMenu(contextMenu, trigger)
        const updatePosition = () => positionContextMenu(contextMenu, trigger)
        window.addEventListener('scroll', updatePosition, true)
        window.addEventListener('resize', updatePosition)
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

  // ------------------------------------------------------------------
  // Table column resizing
  // ------------------------------------------------------------------

  // Adds a drag handle to each resizable header cell except the elastic last one.
  function attachTableResizers() {
    document.querySelectorAll('.popui-table-resizable').forEach(function (table) {
      table.querySelectorAll('thead th').forEach(function (th, idx, all) {
        if (idx === all.length - 1) return
        if (th.querySelector('.popui-table-resizer')) return
        const handle = document.createElement('div')
        handle.className = 'popui-table-resizer'
        th.appendChild(handle)
      })
    })
  }

  // Resizes a column by dragging its handle, freezing the other columns so only it and the elastic last column change size.
  document.addEventListener('mousedown', function (e) {
    const handle = e.target.closest('.popui-table-resizer')
    if (!handle) return
    e.preventDefault()
    const th = handle.closest('th')
    if (!th) return
    const cells = Array.prototype.slice.call(th.parentElement.children)
    const widths = cells.map(function (cell) {
      return cell.offsetWidth
    })
    cells.forEach(function (cell, i) {
      if (i === cells.length - 1) {
        if (!cell.style.minWidth) cell.style.minWidth = widths[i] + 'px'
      } else if (!cell.style.width) {
        cell.style.width = widths[i] + 'px'
        cell.style.minWidth = widths[i] + 'px'
      }
    })
    const startX = e.clientX
    const startWidth = widths[cells.indexOf(th)]
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
    function onMove(mv) {
      const next = Math.max(60, startWidth + (mv.clientX - startX))
      th.style.width = next + 'px'
      th.style.minWidth = next + 'px'
    }
    function onUp() {
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
    }
    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseup', onUp)
  })

  // ------------------------------------------------------------------
  // Card deck reordering
  // ------------------------------------------------------------------

  // The deck's own child cards, in DOM order. Only direct children count, so
  // a Card nested inside another Card never becomes a sortable item.
  function cardDeckCards(root) {
    return Array.prototype.filter.call(root.children, (el) => el.matches && el.matches('[data-card]'))
  }

  // Walks up from an event target to the card that is a direct child of the deck.
  function cardDeckOwnCard(root, target) {
    let card = target && target.closest ? target.closest('[data-card]') : null
    while (card && card.parentElement && card.parentElement !== root) {
      card = card.parentElement.closest('[data-card]')
    }
    return card && card.parentElement === root ? card : null
  }

  // A card's declared position, or Infinity so cards without one sink to the
  // end of the deck rather than jumping to the front.
  function cardDeckOrderOf(card) {
    const v = parseFloat(card.getAttribute('data-order'))
    return isNaN(v) ? Infinity : v
  }

  // Holds the deck's one invariant: data-order, DOM sequence and visible
  // sequence are the same list, numbered 1..n. Sorting is only needed on load,
  // when the server's data-order may disagree with the rendered sequence; a
  // committed move has already put the DOM right and just needs renumbering.
  function cardDeckApplyOrder(root, sort) {
    let cards = cardDeckCards(root)
    if (sort && cards.length) {
      // The original index breaks ties, so equal (and absent) orders keep the
      // sequence they were rendered in.
      const sorted = cards
        .map((card, i) => ({ card, i }))
        .sort((a, b) => (cardDeckOrderOf(a.card) - cardDeckOrderOf(b.card)) || a.i - b.i)
        .map((entry) => entry.card)
      if (sorted.some((card, i) => card !== cards[i])) {
        // Re-inserted before whatever followed the last card, so the head and
        // the hidden order input keep their place.
        const anchor = cards[cards.length - 1].nextSibling
        sorted.forEach((card) => root.insertBefore(card, anchor))
        cards = sorted
      }
    }
    cards.forEach((card, i) => card.setAttribute('data-order', String(i + 1)))
    return cards
  }

  // Turns the deck's buttons off while reordering — a card's menu or action
  // shouldn't be reachable, and it should look unreachable — sparing the
  // reorder toggle itself. Only buttons this disabled are re-enabled, so a
  // button the consumer had already disabled stays that way.
  function cardDeckSetControlsDisabled(root, disabled) {
    root.querySelectorAll('button').forEach((btn) => {
      if (btn.hasAttribute('data-card-deck-reorder')) return
      if (disabled) {
        if (btn.disabled) return
        btn.disabled = true
        btn.setAttribute('data-card-deck-reorder-disabled', '')
      } else if (btn.hasAttribute('data-card-deck-reorder-disabled')) {
        btn.disabled = false
        btn.removeAttribute('data-card-deck-reorder-disabled')
      }
    })
  }

  // Fallback identities for cards carrying neither a data-card-id nor an
  // element id, keyed by the element so they outlive any reordering.
  const cardDeckFallbackIds = new WeakMap()

  // The cards' reported ids, in the order given. A card without an explicit
  // identity gets its 1-based position the first time it is seen, and keeps it
  // from then on: re-deriving the id from the live index would relabel every
  // card on every move, so the reported order would come back identical no
  // matter how the deck was rearranged and the permutation would be lost.
  function cardDeckOrderIds(cards) {
    return cards.map((card, i) => {
      const explicit = card.getAttribute('data-card-id') || card.id
      if (explicit) return explicit
      if (!cardDeckFallbackIds.has(card)) cardDeckFallbackIds.set(card, String(i + 1))
      return cardDeckFallbackIds.get(card)
    })
  }

  // Drags a card vertically within its deck. The dragged card follows the
  // pointer while the cards it passes slide into the slot it left, and on
  // release it animates into its new slot before the DOM is reordered — the
  // card is already sitting where it will land, so the commit is invisible.
  // Calls commit() only when the order actually changed.
  function cardDeckStartDrag(root, card, event, commit) {
    const cards = cardDeckCards(root)
    const from = cards.indexOf(card)
    if (from < 0) return
    // Positions are measured once: every later position is derived from them,
    // so a card sliding out of the way can't feed back into the maths.
    const boxes = cards.map((c) => ({ top: c.offsetTop, height: c.offsetHeight }))
    const centers = boxes.map((b) => b.top + b.height / 2)
    const gap = parseFloat(getComputedStyle(root).rowGap) || 0
    // Removing the dragged card from the flow closes a gap this tall.
    const slot = boxes[from].height + gap
    const startY = event.clientY
    let to = from

    root.classList.add('popui-card-deck-dragging')
    card.classList.add('popui-card-deck-card-dragging')
    document.body.style.cursor = 'grabbing'
    document.body.style.userSelect = 'none'
    // Capture keeps the moves coming even when the pointer outruns the card.
    // Browsers reject ids they no longer track, which is harmless here.
    try { card.setPointerCapture(event.pointerId) } catch (err) { /* no capture */ }

    // Shifts every other card by one slot when it sits between the card's old
    // and new index — the exact offset it gets once the drag is committed.
    function layout() {
      cards.forEach((c, i) => {
        if (i === from) return
        let dy = 0
        if (to > from && i > from && i <= to) dy = -slot
        else if (to < from && i >= to && i < from) dy = slot
        c.style.transform = dy ? 'translateY(' + dy + 'px)' : ''
      })
    }

    function onMove(e) {
      const dy = e.clientY - startY
      // A swap happens as the dragged card's leading edge crosses the centre
      // of its neighbour, which is where the exchange reads as complete.
      const top = boxes[from].top + dy
      const bottom = top + boxes[from].height
      let next = from
      while (next > 0 && top < centers[next - 1]) next--
      while (next < cards.length - 1 && bottom > centers[next + 1]) next++
      if (next !== to) {
        to = next
        layout()
      }
      card.style.transform = 'translateY(' + dy + 'px)'
    }

    function onUp() {
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerup', onUp)
      window.removeEventListener('pointercancel', onUp)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
      // Where the card comes to rest: the slot it now owns, measured against
      // the layout captured at drag start.
      let rest = 0
      if (to > from) rest = boxes[to].top + boxes[to].height - boxes[from].height - boxes[from].top
      else if (to < from) rest = boxes[to].top - boxes[from].top
      const current = card.style.transform
      const target = rest ? 'translateY(' + rest + 'px)' : ''
      // Dropping the dragging class restores the transform transition, so the
      // card glides the last few pixels into place.
      card.classList.remove('popui-card-deck-card-dragging')
      card.style.transform = target

      let settled = false
      function settle() {
        if (settled) return
        settled = true
        // Committing is a swap of two equal quantities: every card's layout
        // position shifts by exactly the offset it was holding, so the reorder
        // and the transform reset cancel out to no visible movement — as long
        // as they land in one frame with transitions off. Animating between
        // them instead makes each displaced card overshoot by a slot and slide
        // back, even though it was already sitting in the right place.
        root.classList.add('popui-card-deck-settling')
        if (to !== from) {
          if (to > from) cards[to].after(card)
          else cards[to].before(card)
        }
        cards.forEach((c) => { c.style.transform = '' })
        // Flushes layout so the untransitioned state is what the frame paints.
        void root.offsetHeight
        root.classList.remove('popui-card-deck-settling')
        root.classList.remove('popui-card-deck-dragging')
        if (to !== from) commit()
      }
      if (current === target) settle()
      else {
        card.addEventListener('transitionend', settle, { once: true })
        setTimeout(settle, 300)
      }
    }

    window.addEventListener('pointermove', onMove)
    window.addEventListener('pointerup', onUp)
    window.addEventListener('pointercancel', onUp)
  }


  // ------------------------------------------------------------------
  // Field selector helpers
  // ------------------------------------------------------------------

  // Splits text into matched and unmatched runs, so a field row can highlight
  // the part of its path a query hit. Overlapping matches are merged rather
  // than split twice.
  function highlightParts(text, terms) {
    const lower = text.toLowerCase()
    const ranges = []

    for (const term of terms) {
      const needle = term.toLowerCase()
      let from = lower.indexOf(needle)
      while (from !== -1) {
        ranges.push([from, from + needle.length])
        from = lower.indexOf(needle, from + needle.length)
      }
    }

    if (!ranges.length) return [{ m: false, t: text }]

    ranges.sort((a, b) => a[0] - b[0])

    const parts = []
    let cursor = 0
    for (const [start, end] of ranges) {
      if (end <= cursor) continue
      const from = Math.max(start, cursor)
      if (from > cursor) parts.push({ m: false, t: text.slice(cursor, from) })
      parts.push({ m: true, t: text.slice(from, end) })
      cursor = end
    }
    if (cursor < text.length) parts.push({ m: false, t: text.slice(cursor) })

    return parts
  }

  // Scores one field against the query's terms, or null when any term misses.
  // Only the path is consulted: the description is prose, and matching it
  // would surface fields for words that merely occur in their explanation.
  function scoreField(entry, terms) {
    const path = entry.lower
    const name = entry.name.toLowerCase()

    // Shallow fields win ties: `code` should beat `lines[].item.code`.
    let score = -entry.depth * 3

    for (const term of terms) {
      const index = path.indexOf(term)
      if (index >= 0) {
        score += 100 - Math.min(index, 60)
        if (name === term) score += 80
        else if (name.startsWith(term)) score += 40
        continue
      }
      if (isSubsequence(term, path)) {
        score += 20
        continue
      }
      return null
    }

    return score
  }

  const TYPE_ABBREVIATIONS = {
    object: 'obj',
    string: 'str',
    array: 'arr',
    integer: 'int',
    number: 'num',
    boolean: 'bool',
  }

  // Escapes text for x-html. Field names come from schemas, but a tree can
  // be hand-written too.
  function escapeHTML(text) {
    return String(text).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  }

  // Reports whether the characters of term appear in text in order, which is
  // what lets `supname` find `supplier.name`.
  function isSubsequence(term, text) {
    let cursor = 0
    for (const character of text) {
      if (character === term[cursor]) cursor++
      if (cursor === term.length) return true
    }
    return false
  }

  // ------------------------------------------------------------------
  // Template editor helpers (Contenteditable with a VariableFormat)
  // ------------------------------------------------------------------

  // Builds the regular expression that finds variables written in a format
  // such as "{{.%s}}": everything in the format is literal except %s, which
  // stands for the name.
  function variablePattern(format, whole) {
    const escaped = String(format).replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    const source = escaped.replace('%s', '(.+?)')
    // The whole-string form is used with match(), which only reports the
    // capture group — the variable's name — without the global flag.
    return whole ? new RegExp('^' + source + '$') : new RegExp(source, 'g')
  }

  // The chip a variable is drawn as. data-variable carries the variable as
  // written, so serialising the editor gives the template back verbatim.
  function variableChipHTML(variable, name) {
    return '<span class="tag" contenteditable="false" data-variable="' + escapeHTML(variable) + '">' + escapeHTML(name) + '</span>'
  }

  function variableChipNode(variable, name) {
    const span = document.createElement('span')
    span.className = 'tag'
    span.setAttribute('contenteditable', 'false')
    span.setAttribute('data-variable', variable)
    span.textContent = name
    return span
  }

  // Template text → editor HTML: variables become chips, newlines line breaks.
  function renderTemplate(value, format) {
    const pattern = variablePattern(format, false)
    let html = ''
    let last = 0
    for (const m of String(value).matchAll(pattern)) {
      html += escapeHTML(value.slice(last, m.index))
      html += variableChipHTML(m[0], m[1])
      last = m.index + m[0].length
    }
    html += escapeHTML(value.slice(last))
    return html.replace(/\n/g, '<br>')
  }

  // Editor DOM → template text: chips give back the variable they carry,
  // <br> and the <div>/<p> blocks browsers create on Enter become newlines,
  // and the non-breaking spaces browsers put in contenteditables become
  // ordinary ones.
  function serializeTemplate(root) {
    let out = ''
    const walk = (node) => {
      for (const child of node.childNodes) {
        if (child.nodeType === 3) {
          out += child.nodeValue.replace(/\u00a0/g, ' ')
          continue
        }
        if (child.nodeType !== 1) continue
        if (child.hasAttribute('data-variable')) {
          out += child.getAttribute('data-variable')
          continue
        }
        const tag = child.tagName
        if (tag === 'BR') {
          out += '\n'
          continue
        }
        if ((tag === 'DIV' || tag === 'P') && out && !out.endsWith('\n')) out += '\n'
        walk(child)
      }
    }
    walk(root)
    return out
  }

  // ------------------------------------------------------------------
  // Alpine controllers
  // ------------------------------------------------------------------

  document.addEventListener('alpine:init', () => {
    if (!window.Alpine) return

    // Inline option list for filter chips: owns the selection, the arrow-key highlight, and the open state of its panel.
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
        // Keys are handled at the document level so they keep working when a closing popover moves focus.
        this._onKeydown = (e) => {
          if (!this.$root || this.$root.offsetParent === null) return
          if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp' && e.key !== 'Enter' && e.key !== 'Escape' && e.key !== ' ') return
          const a = document.activeElement
          const editable = a && (a.tagName === 'INPUT' || a.tagName === 'TEXTAREA' || a.tagName === 'SELECT' || a.isContentEditable)
          if (editable && !this.$root.contains(a)) return
          if (!this.open) {
            if (this.$root.contains(a) && (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ')) {
              e.preventDefault()
              this.openPanel()
            }
            return
          }
          const form = this.$root.closest('form')
          if (a && a !== document.body && !this.$root.contains(a) && !(form && form.contains(a))) return
          this.onKeydown(e)
        }
        document.addEventListener('keydown', this._onKeydown, true)
        // A click outside the chip closes the options panel.
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
        submitClosestForm(this.$root)
      },
      openPanel() {
        this.open = true
        if (this.activeIndex < 0 && this.optionValues.length) this.activeIndex = 0
        this.$nextTick(() => clampPanelX(this.$refs.panel))
      },
      closePanel() {
        this.open = false
        this.activeIndex = -1
      },
      togglePanel() {
        this.open ? this.closePanel() : this.openPanel()
      },
      // Moves the highlight by delta, wrapping around the option list.
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
          e.preventDefault()
          if (this.activeIndex >= 0) this.toggle(this.optionValues[this.activeIndex])
        } else if (e.key === 'Escape') {
          this.closePanel()
        }
      },
      // Highlights a row and toggles it, keeping mouse clicks and Enter in sync.
      choose(i) {
        if (i < 0 || i >= this.optionValues.length) return
        this.activeIndex = i
        this.toggle(this.optionValues[i])
      },
      // Applies a value and submits after Alpine has rendered the hidden inputs.
      toggle(v) {
        this.values = toggleValue(this.values, v, this.multiple)
        this.initial = JSON.stringify(this.values)
        this.$nextTick(() => {
          clampPanelX(this.$refs.panel)
          this.submitForm()
        })
      },
    }))

    // Filter row: tracks the active filter chips, adds and removes them, and resets their editors.
    Alpine.data('filterRow', (initialActive, allNames) => ({
      active: Array.isArray(initialActive) ? [...initialActive] : [],
      all: Array.isArray(allNames) ? [...allNames] : [],

      // Returns the enclosing form, falling back to the controller root.
      _form() {
        return (this.$root.closest && this.$root.closest('form')) || this.$root
      },
      // Submits the form after Alpine has removed the cleared hidden inputs.
      _submit() {
        this.$nextTick(() => submitClosestForm(this._form()))
      },
      // Resets every editor a chip may contain: text input, select, dropdown, option list, or calendar.
      clearFilter(name) {
        const chip = this._form().querySelector('[data-filter-name="' + name + '"]')
        if (!chip) return
        const input = chip.querySelector('input[type="text"][name="' + name + '"]')
        if (input) input.value = ''
        const select = chip.querySelector('select[name="' + name + '"]')
        if (select) select.value = ''
        const dropdown = alpineData(chip.querySelector('[role="combobox"]'))
        if (dropdown && Array.isArray(dropdown.values)) {
          dropdown.values = []
          dropdown.initial = '[]'
        }
        const optionList = alpineData(chip.querySelector('[data-filter-options]'))
        if (optionList && Array.isArray(optionList.values)) {
          optionList.values = []
          optionList.initial = '[]'
          optionList.activeIndex = -1
        }
        const calendar = alpineData(chip.querySelector('[data-filter-calendar]'))
        if (calendar && typeof calendar.clear === 'function') calendar.clear()
      },
      // Opens a chip's value editor right after it appears.
      autoOpenChip(name) {
        const chip = this._form().querySelector('[data-filter-name="' + name + '"]')
        if (!chip) return
        const optionList = chip.querySelector('[data-filter-options]')
        if (optionList) {
          const d = alpineData(optionList)
          if (d && typeof d.openPanel === 'function') d.openPanel()
          // Focus is asserted twice because a closing popover can restore focus to its invoker a beat later.
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
        const calendar = chip.querySelector('[data-filter-calendar]')
        if (calendar) {
          const d = alpineData(calendar)
          if (d && typeof d.openPanel === 'function') d.openPanel()
          if (typeof calendar.focus === 'function') calendar.focus()
          return
        }
        // Popovers are retried once because browsers serialise popover transitions.
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
      // Appends a filter to the active list and opens its editor.
      add(name) {
        if (!this.active.includes(name)) this.active = [...this.active, name]
        this.$nextTick(() => this.autoOpenChip(name))
      },
      // Returns a chip's flex order so chips lay out in the order they were added.
      orderOf(name) {
        const i = this.active.indexOf(name)
        return i < 0 ? 0 : i
      },
      // Removes one filter, resets its editor, and submits.
      remove(name) {
        this.active = this.active.filter((n) => n !== name)
        this.clearFilter(name)
        this._submit()
      },
      // Removes every active filter, resets their editors, and submits once.
      clearAll() {
        const names = [...this.active]
        this.active = []
        names.forEach((n) => this.clearFilter(n))
        this._submit()
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
      // Reports whether any field is still inactive, which keeps the add button visible.
      hasAvailable() {
        return this.all.some((n) => !this.active.includes(n))
      },
    }))


    // Template editor: a Contenteditable that draws each variable as a chip and keeps the template text —
    // variables written out — as its value. Rich view edits the chips in place; plain view edits the text.
    Alpine.data('templateEditor', (init) => ({
      value: (init && init.value) || '',
      format: (init && init.format) || '{{.%s}}',
      view: (init && init.view) || 'rich',
      savedRange: null,

      init() {
        this.render()
        // The chips are rebuilt from the value whenever the rich view comes
        // back, since the plain view may have changed the text.
        this.$watch('view', (view) => {
          // The views trade places at the height the last one had, so the
          // box does not jump: the rich one grows with its content while the
          // textarea keeps whatever height it was given.
          const from = view === 'rich' ? this.$refs.plain : this.$refs.editor
          const to = view === 'rich' ? this.$refs.editor : this.$refs.plain
          if (from && to && from.offsetHeight) to.style.minHeight = from.offsetHeight + 'px'
          if (view === 'rich') this.$nextTick(() => this.render())
        })
        // Something else may set the value — an x-model from outside, a
        // fetch — so a change that did not come from the editor redraws it.
        this.$watch('value', (value) => { if (this.view === 'rich' && serializeTemplate(this.$refs.editor) !== value) this.render() })
        // The caret is tracked while the editor has the selection: a picker
        // that inserts into it takes focus first, and the insert has to go
        // where the caret was.
        this._onSelectionChange = () => {
          const sel = window.getSelection()
          if (!sel || !sel.rangeCount) return
          const range = sel.getRangeAt(0)
          if (this.$refs.editor && this.$refs.editor.contains(range.commonAncestorContainer)) this.savedRange = range.cloneRange()
        }
        document.addEventListener('selectionchange', this._onSelectionChange)
      },
      destroy() {
        if (this._onSelectionChange) document.removeEventListener('selectionchange', this._onSelectionChange)
      },

      render() {
        if (this.$refs.editor) this.$refs.editor.innerHTML = renderTemplate(this.value, this.format)
      },
      // Called on every edit of the rich view. A variable typed out by hand
      // — the closing brace just landed — is turned into a chip right away,
      // which means redrawing the editor; the caret is carried across as an
      // offset into the template text, which both the old and the new DOM
      // serialise to.
      sync() {
        const editor = this.$refs.editor
        this.value = serializeTemplate(editor)
        if (this.hasTypedVariable(editor)) {
          const offset = this.caretOffset(editor)
          this.render()
          if (offset !== null) this.setCaretOffset(editor, offset)
        }
        this.changed()
      },
      // Announces a change with an input event from the component root —
      // after Alpine has carried the new value across an x-model on it, so
      // a handler on the root reads page state that is already current.
      // The views' own input events stop at them for the same reason: let
      // through, they would reach the root a tick too early.
      changed() {
        this.$nextTick(() => this.$root.dispatchEvent(new Event('input', { bubbles: true })))
      },
      // Reports whether any text in the editor — outside the chips — is a
      // complete variable.
      hasTypedVariable(editor) {
        editor.normalize()
        const pattern = variablePattern(this.format, false)
        const walker = document.createTreeWalker(editor, NodeFilter.SHOW_TEXT)
        for (let node = walker.nextNode(); node; node = walker.nextNode()) {
          if (node.parentElement && node.parentElement.closest('[data-variable]')) continue
          pattern.lastIndex = 0
          if (pattern.test(node.nodeValue)) return true
        }
        return false
      },
      // The caret's position as a number of characters into the template.
      caretOffset(editor) {
        const sel = window.getSelection()
        if (!sel || !sel.rangeCount) return null
        const range = sel.getRangeAt(0)
        if (!editor.contains(range.startContainer)) return null
        const before = document.createRange()
        before.setStart(editor, 0)
        before.setEnd(range.startContainer, range.startOffset)
        return serializeTemplate(before.cloneContents()).length
      },
      // Puts the caret at a template offset in a freshly rendered editor,
      // whose children are only text, chips and line breaks.
      setCaretOffset(editor, offset) {
        const place = (node, at) => {
          const range = document.createRange()
          range.setStart(node, at)
          range.collapse(true)
          const sel = window.getSelection()
          sel.removeAllRanges()
          sel.addRange(range)
          this.savedRange = range.cloneRange()
        }
        let remaining = offset
        const children = Array.from(editor.childNodes)
        for (let i = 0; i < children.length; i++) {
          const child = children[i]
          if (child.nodeType === 3) {
            if (remaining <= child.nodeValue.length) return place(child, remaining)
            remaining -= child.nodeValue.length
            continue
          }
          if (child.nodeType !== 1) continue
          const length = child.hasAttribute('data-variable') ? child.getAttribute('data-variable').length : child.tagName === 'BR' ? 1 : 0
          if (remaining === 0) return place(editor, i)
          if (remaining <= length) return place(editor, i + 1)
          remaining -= length
        }
        place(editor, children.length)
      },
      isVariable(text) {
        return variablePattern(this.format, true).test(text)
      },
      // Puts text at the caret of whichever view is showing. In the rich view
      // a variable becomes a chip and anything else plain text; the caret
      // ends up after it. This is what FieldPicker calls.
      insertText(text) {
        if (this.view === 'plain') {
          const area = this.$refs.plain
          if (!area) return
          const start = area.selectionStart == null ? area.value.length : area.selectionStart
          const end = area.selectionEnd == null ? start : area.selectionEnd
          area.setRangeText(text, start, end, 'end')
          area.focus()
          this.value = area.value
          this.changed()
          return
        }
        const editor = this.$refs.editor
        editor.focus()
        const sel = window.getSelection()
        let range = this.savedRange && editor.contains(this.savedRange.commonAncestorContainer) ? this.savedRange : null
        if (!range) {
          range = document.createRange()
          range.selectNodeContents(editor)
          range.collapse(false)
        }
        range.deleteContents()
        const match = String(text).match(variablePattern(this.format, true))
        const node = match ? variableChipNode(text, match[1]) : document.createTextNode(text)
        range.insertNode(node)
        range.setStartAfter(node)
        range.collapse(true)
        sel.removeAllRanges()
        sel.addRange(range)
        this.savedRange = range.cloneRange()
        this.sync()
      },
    }))

    // Field picker over a nested data shape (usually a GOBL document): filters the whole tree by path, or browses it a level at a time.
    // The tree arrives as JSON in a script element; paths are rebuilt here from each entry's parent index. Picking a field
    // formats its path, inserts it at the caret of the target element when there is one, and raises field-select.
    Alpine.data('fieldPicker', (init) => {
      // The entries live here as plain objects as well as on the component,
      // so the scorer reads them without going through Alpine's proxies —
      // two thousand fields times a few properties per keystroke adds up.
      const list = []
      // rows is a getter, and the template asks for it once per row (every
      // row's highlight compares against activeIndex, which reads rows), so
      // without a cache one keystroke would rank the whole tree eighty times.
      let rowsKey = null
      let rowsCache = []
      const isTextInput = (el) => !!el && (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA')

      return {
        entries: [],
        roots: [],
        value: (init && init.value) || '',
        name: (init && init.name) || '',
        // What a picked path is wrapped in: %s stands for the path. A field
        // with a Value of its own skips it.
        format: (init && init.format) || '%s',
        // Selector of the element a pick is inserted into, if any.
        target: (init && init.target) || '',
        disabled: !!(init && init.disabled),
        // Browse-only objects and arrays: they open but are not picked.
        scalarsOnly: !!(init && init.scalarsOnly),
        query: '',
        open: false,
        // The objects drilled into, root first, as entry indices.
        trail: [],
        // The keyboard highlight follows a path rather than a row index, so it
        // survives the list changing underneath it.
        activePath: '',
        // The last caret position seen inside a contenteditable target. The
        // filter box takes focus when the panel opens, which moves the
        // document selection away from the target, so it is tracked while
        // the target still has it.
        savedRange: null,

        init() {
          const el = init && init.source ? document.getElementById(init.source) : null
          let raw = []
          if (el) {
            try { raw = JSON.parse(el.textContent) || [] } catch (e) { raw = [] }
          } else if (init && Array.isArray(init.fields)) {
            raw = init.fields
          }
          list.length = 0
          raw.forEach((e) => list.push({
            name: e.n || '',
            type: e.y || '',
            description: e.d || '',
            own: e.v || '',
            parent: typeof e.p === 'number' ? e.p : -1,
            array: !!(e.f & 1),
            always: !!(e.f & 2),
            children: !!(e.f & 4),
            path: '',
            lower: '',
            depth: 0,
            kids: [],
          }))
          list.forEach((entry, i) => {
            const parent = entry.parent >= 0 ? list[entry.parent] : null
            entry.path = (parent ? parent.path + '.' : '') + entry.name + (entry.array ? '[]' : '')
            entry.lower = entry.path.toLowerCase()
            entry.depth = parent ? parent.depth + 1 : 1
            if (parent) parent.kids.push(i)
            else this.roots.push(i)
          })
          this.entries = list
          rowsKey = null

          if (this.target) {
            this._onSelectionChange = () => {
              const target = this.targetEl()
              if (!target || isTextInput(target)) return
              const sel = window.getSelection()
              if (!sel || !sel.rangeCount) return
              const range = sel.getRangeAt(0)
              if (target.contains(range.commonAncestorContainer)) this.savedRange = range.cloneRange()
            }
            document.addEventListener('selectionchange', this._onSelectionChange)
          }
        },
        destroy() {
          if (this._onSelectionChange) document.removeEventListener('selectionchange', this._onSelectionChange)
        },

        get searching() {
          return this.query.trim().length > 0
        },
        // The entries of the level currently browsed.
        get level() {
          if (!this.trail.length) return this.roots
          return this.entries[this.trail[this.trail.length - 1]].kids
        },
        get rows() {
          const key = this.query + '\n' + this.trail.join(',')
          if (key !== rowsKey) {
            rowsKey = key
            rowsCache = this.searching ? this.match(this.query) : this.level.filter((i) => this.shown(i)).map((i) => ({ i, parts: null }))
          }
          return rowsCache
        },
        // With scalarsOnly, a non-scalar with nothing underneath — an array
        // of strings, a free-form map — can be neither picked nor opened, so
        // it is left out rather than shown disabled.
        shown(i) {
          return this.pickable(i) || this.entries[i].kids.length > 0
        },
        get activeIndex() {
          const i = this.rows.findIndex((row) => this.entries[row.i].path === this.activePath)
          return i === -1 ? 0 : i
        },
        get activeId() {
          if (!this.open || !this.rows.length) return null
          return this.$id('field-picker') + '-opt-' + this.activeIndex
        },

        // The value a field emits: its own override when it has one, otherwise
        // its path wrapped in the format. Not named valueOf: Alpine resolves an
        // expression's names against a proxy that inherits Object.prototype,
        // so a method by that name is shadowed and silently returns the scope.
        emitted(i) {
          const entry = this.entries[i]
          if (!entry) return ''
          if (entry.own) return entry.own
          return this.format.includes('%s') ? this.format.replace('%s', entry.path) : this.format + entry.path
        },
        isSelected(i) {
          return !!this.value && this.emitted(i) === this.value
        },
        // Whether picking the field is allowed at all: with scalarsOnly an
        // object or array is there to be opened, not taken as a value.
        pickable(i) {
          const entry = this.entries[i]
          if (!entry) return false
          return !this.scalarsOnly || (entry.type !== 'object' && entry.type !== 'array')
        },
        typeLabel(i) {
          const type = this.entries[i] ? this.entries[i].type : ''
          return TYPE_ABBREVIATIONS[type] || type
        },
        // The row's path as escaped HTML: the whole path while browsing, and
        // while filtering the runs a term did not hit dimmed, so the ones it
        // did still read at full strength.
        labelHTML(row) {
          const entry = this.entries[row.i]
          if (!row.parts) return escapeHTML(entry.name + (entry.array ? '[]' : ''))
          let html = ''
          for (const part of row.parts) {
            html += part.m ? escapeHTML(part.t) : '<span class="text-foreground-default-secondary">' + escapeHTML(part.t) + '</span>'
          }
          return html
        },

        // Ranks every field against the query. Terms are matched against the
        // whole path, so `sup name` and `supplier.name` find the same field.
        match(query, limit = 80) {
          const terms = query.toLowerCase().split(/\s+/).filter(Boolean)
          if (!terms.length) return []

          const found = []
          for (let i = 0; i < list.length; i++) {
            if (!this.shown(i)) continue
            const score = scoreField(list[i], terms)
            if (score === null) continue
            found.push({ i, score })
          }

          found.sort((a, b) => {
            const x = list[a.i]
            const y = list[b.i]
            return b.score - a.score || x.path.length - y.path.length || (x.path < y.path ? -1 : x.path > y.path ? 1 : 0)
          })

          // Only the rows that are rendered are worth highlighting.
          return found.slice(0, limit).map((row) => ({ i: row.i, parts: highlightParts(list[row.i].path, terms) }))
        },

        move(delta) {
          if (!this.rows.length) return
          const next = (this.activeIndex + delta + this.rows.length) % this.rows.length
          this.activePath = this.entries[this.rows[next].i].path
          this.scrollActiveIntoView()
        },
        drill(i) {
          const entry = this.entries[i]
          if (!entry || !entry.kids.length) return
          this.trail = [...this.trail, i]
          this.query = ''
          this.activePath = ''
        },
        back() {
          if (!this.trail.length) return
          const parent = this.trail[this.trail.length - 1]
          this.trail = this.trail.slice(0, -1)
          this.activePath = this.entries[parent].path
        },
        select(i) {
          const entry = this.entries[i]
          if (!entry) return
          // A browse-only row's only action is to open it.
          if (!this.pickable(i)) {
            if (entry.kids.length) this.drill(i)
            return
          }
          this.value = this.emitted(i)
          const target = this.targetEl()
          if (target) this.insert(this.value, target)
          this.$root.dispatchEvent(new CustomEvent('field-select', {
            bubbles: true,
            detail: { path: entry.path, value: this.value, type: entry.type, inserted: !!target },
          }))
          this.$refs.panel.hidePopover()
        },
        targetEl() {
          if (!this.target) return null
          try { return document.querySelector(this.target) } catch (e) { return null }
        },
        // Puts text at the caret of an input, textarea or contenteditable
        // element, replacing any selection there, and leaves the caret after
        // it. The caret is where it was when the target last had focus:
        // inputs remember that themselves, a contenteditable is tracked via
        // savedRange, and with nothing to go on the text goes at the end.
        insert(text, target = this.targetEl()) {
          if (!target) return
          // A template editor draws variables as chips, so it does the
          // inserting itself.
          if (target.hasAttribute('data-template-editor') && window.Alpine) {
            const editor = window.Alpine.$data(target)
            if (editor && typeof editor.insertText === 'function') {
              editor.insertText(text)
              return
            }
          }
          if (isTextInput(target)) {
            const start = target.selectionStart == null ? target.value.length : target.selectionStart
            const end = target.selectionEnd == null ? start : target.selectionEnd
            target.setRangeText(text, start, end, 'end')
            target.focus()
            target.dispatchEvent(new Event('input', { bubbles: true }))
            return
          }
          target.focus()
          const sel = window.getSelection()
          let range = this.savedRange && target.contains(this.savedRange.commonAncestorContainer) ? this.savedRange : null
          if (!range) {
            range = document.createRange()
            range.selectNodeContents(target)
            range.collapse(false)
          }
          sel.removeAllRanges()
          sel.addRange(range)
          // insertText keeps the edit on the undo stack and fires the input
          // events itself; the manual path is for engines without it.
          let done = false
          try { done = document.execCommand('insertText', false, text) } catch (e) { done = false }
          if (!done) {
            range.deleteContents()
            const node = document.createTextNode(text)
            range.insertNode(node)
            range.setStartAfter(node)
            range.collapse(true)
            sel.removeAllRanges()
            sel.addRange(range)
            target.dispatchEvent(new Event('input', { bubbles: true }))
          }
          if (sel.rangeCount) this.savedRange = sel.getRangeAt(0).cloneRange()
        },
        scrollActiveIntoView() {
          this.$nextTick(() => {
            const el = document.getElementById(this.activeId)
            if (el) el.scrollIntoView({ block: 'nearest' })
          })
        },
        // Opens where the current value lives rather than at the root.
        openAtValue() {
          this.trail = []
          this.activePath = ''
          if (!this.value) return
          const found = this.entries.findIndex((_, i) => this.isSelected(i))
          if (found === -1) return
          this.activePath = this.entries[found].path
          const trail = []
          for (let i = this.entries[found].parent; i >= 0; i = this.entries[i].parent) trail.unshift(i)
          this.trail = trail
        },
        show() {
          if (this.disabled || this.open) return
          this.$refs.panel.showPopover()
        },

        onKeydown(e) {
          // Buttons inside the panel — the crumbs, Clear — keep their own Enter.
          if (e.key === 'Enter' && e.target && e.target.tagName === 'BUTTON' && this.$refs.panel.contains(e.target)) return
          if (e.key === 'ArrowDown') {
            e.preventDefault()
            if (!this.open) this.show()
            else this.move(1)
          } else if (e.key === 'ArrowUp') {
            e.preventDefault()
            if (!this.open) this.show()
            else this.move(-1)
          } else if (e.key === 'Enter') {
            if (!this.open) return
            e.preventDefault()
            const row = this.rows[this.activeIndex]
            if (row) this.select(row.i)
          } else if (e.key === 'Escape') {
            if (!this.open) return
            e.preventDefault()
            // Escape clears the filter first, so a mistyped query does not cost
            // the whole dropdown.
            if (this.query) this.query = ''
            else this.$refs.panel.hidePopover()
          } else if (this.open && !this.searching && (e.key === 'ArrowRight' || e.key === 'ArrowLeft')) {
            // Left and Right move the caret while there is something to filter.
            e.preventDefault()
            const row = this.rows[this.activeIndex]
            if (e.key === 'ArrowRight') { if (row) this.drill(row.i) } else this.back()
          }
        },
        onToggle(e) {
          if (e.newState === 'open') {
            this.open = true
            // At least as wide as the trigger, so a full-width field gets a
            // panel to match; the stylesheet sets the floor for a button.
            const trigger = this.$refs.field.firstElementChild || this.$refs.field
            this.$refs.panel.style.minWidth = trigger.offsetWidth + 'px'
            this.query = ''
            this.openAtValue()
            this.$nextTick(() => { if (this.$refs.search) this.$refs.search.focus() })
          } else {
            this.open = false
            this.query = ''
          }
        },
      }
    })

    // Dual-month date-range picker with a preset rail and a Cancel / Confirm footer; only Confirm applies the pending selection.
    // With single: true it becomes a single-date picker: one month grid, no presets, and a day click sets from = to.
    Alpine.data('rangeCalendar', (init) => ({
      name: (init && init.name) || '',
      single: !!(init && init.single),
      open: false,
      preset: 'custom',
      // Pending selection edited by the grids.
      from: (init && init.from) || null,
      to: (init && init.to) || null,
      // Committed selection exposed as rangeValue and summary.
      committedFrom: (init && init.from) || null,
      committedTo: (init && init.to) || null,
      viewY: 2000,
      viewM: 0,
      dows: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
      presets: (init && init.presets) || [
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
        // An outside click while the panel is open cancels the selection.
        this._onDocClick = (e) => {
          if (this.open && this.$root && !this.$root.contains(e.target)) this.cancel()
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
      // Builds six weeks of cells for a month, marking outside days and today.
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
        const months = [{ label: this.monthLabel(this.viewY, this.viewM), weeks: this.monthGrid(this.viewY, this.viewM) }]
        if (!this.single) {
          const r = new Date(this.viewY, this.viewM + 1, 1)
          months.push({ label: this.monthLabel(r.getFullYear(), r.getMonth()), weeks: this.monthGrid(r.getFullYear(), r.getMonth()) })
        }
        return months
      },
      prev() {
        const d = new Date(this.viewY, this.viewM - 1, 1)
        this.viewY = d.getFullYear(); this.viewM = d.getMonth()
      },
      next() {
        const d = new Date(this.viewY, this.viewM + 1, 1)
        this.viewY = d.getFullYear(); this.viewM = d.getMonth()
      },
      // Returns a day's range state, with outside-month days always unstyled.
      // A one-day range (from === to) reports 'single' so it rounds on both sides.
      dayState(iso, outside) {
        if (outside) return null
        if (!this.from) return null
        if (!this.to) return iso === this.from ? 'start' : null
        if (iso === this.from) return this.to === this.from ? 'single' : 'start'
        if (iso === this.to) return 'end'
        if (iso > this.from && iso < this.to) return 'middle'
        return null
      },
      // Starts a new range, moves the start when clicking before it, or completes
      // the range. In single mode a click just selects that day (from = to).
      selectDay(iso, outside) {
        if (outside) return
        if (this.single) {
          this.from = iso; this.to = iso
          return
        }
        if (!this.from || (this.from && this.to)) {
          this.from = iso; this.to = null
          this.preset = 'custom'
          return
        }
        if (iso < this.from) { this.from = iso }
        else { this.to = iso }
      },
      // Applies a preset's date range and jumps the view to it.
      setPreset(key) {
        this.preset = key
        if (key === 'custom') {
          const c = new Date(); c.setHours(0, 0, 0, 0)
          this.from = this.iso(c); this.to = null
          return
        }
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
      },
      // Formats an ISO date as dd/mm/yyyy.
      fmt(iso) {
        const p = String(iso).split('-')
        return p[2] + '/' + p[1] + '/' + p[0]
      },
      get summary() {
        if (this.committedFrom && this.committedTo) {
          if (this.single) return this.fmt(this.committedFrom)
          return this.fmt(this.committedFrom) + ' → ' + this.fmt(this.committedTo)
        }
        return ''
      },
      get rangeValue() {
        if (this.committedFrom && this.committedTo) {
          if (this.single) return this.committedFrom
          return this.committedFrom + '..' + this.committedTo
        }
        return ''
      },
      togglePanel() { this.open ? this.open = false : this.openPanel() },
      // Opens the panel and keeps it inside the viewport.
      openPanel() {
        this.open = true
        this.$nextTick(() => clampPanelX(this.$refs.panel))
      },
      // Resets the pending and committed selection without submitting.
      clear() {
        this.from = null; this.to = null
        this.committedFrom = null; this.committedTo = null
        this.preset = 'custom'
      },
      // Commits the pending range, closes the panel, announces it, and submits the enclosing form.
      confirm() {
        if (!this.to) return
        this.committedFrom = this.from
        this.committedTo = this.to
        this.open = false
        if (this.$root) this.$root.dispatchEvent(new CustomEvent('popui-cal-confirm', { bubbles: true }))
        this._submit()
      },
      // Clears the selection, closes the panel, announces it, and resubmits when a committed range was cleared.
      cancel() {
        const hadCommitted = !!(this.committedFrom || this.committedTo)
        this.from = null; this.to = null
        this.committedFrom = null; this.committedTo = null
        this.preset = 'custom'
        this.open = false
        if (this.$root) this.$root.dispatchEvent(new CustomEvent('popui-cal-cancel', { bubbles: true }))
        if (hadCommitted) this._submit()
      },
      _submit() {
        this.$nextTick(() => submitClosestForm(this.$root))
      },
    }))

    // Side panel opened and closed by popui-sidepanel-open/close window events whose detail matches the panel id.
    // The width is user-adjustable by dragging the panel's inner-edge handle.
    Alpine.data('sidePanel', (id, width = 400, anchor = 'right') => ({
      open: false,
      width,
      resizing: false,
      init() {
        const matches = (e) => e && e.detail === id
        this._onOpen = (e) => { if (matches(e)) this.open = true }
        this._onClose = (e) => { if (matches(e)) this.open = false }
        this._onKeydown = (e) => {
          if (e.key === 'Escape' && this.open) {
            this.open = false
            e.stopPropagation()
          }
        }
        window.addEventListener('popui-sidepanel-open', this._onOpen)
        window.addEventListener('popui-sidepanel-close', this._onClose)
        document.addEventListener('keydown', this._onKeydown)
        // Every open-state change is re-broadcast so consumers can react to any close path.
        this.$watch('open', (newVal, oldVal) => {
          if (newVal === oldVal) return
          const name = newVal ? 'popui-sidepanel-open' : 'popui-sidepanel-close'
          window.dispatchEvent(new CustomEvent(name, { detail: id }))
        })
      },
      destroy() {
        window.removeEventListener('popui-sidepanel-open', this._onOpen)
        window.removeEventListener('popui-sidepanel-close', this._onClose)
        document.removeEventListener('keydown', this._onKeydown)
      },
      // Drag-resizes the panel from its inner-edge handle, clamped so the
      // panel stays usable and never covers the whole viewport.
      startResize(e) {
        if (e.button > 0) return
        e.preventDefault()
        this.resizing = true
        const startX = e.clientX
        const startWidth = this.width
        const dir = anchor === 'left' ? 1 : -1
        const min = 320
        const max = Math.max(min, window.innerWidth - 320)
        const onMove = (mv) => {
          this.width = Math.min(max, Math.max(min, startWidth + dir * (mv.clientX - startX)))
        }
        const onUp = () => {
          window.removeEventListener('pointermove', onMove)
          window.removeEventListener('pointerup', onUp)
          document.body.style.cursor = ''
          document.body.style.userSelect = ''
          this.resizing = false
        }
        document.body.style.cursor = 'col-resize'
        document.body.style.userSelect = 'none'
        window.addEventListener('pointermove', onMove)
        window.addEventListener('pointerup', onUp)
      },
    }))

    // Card deck reorder mode: drags the child cards to a new vertical position
    // and reports the resulting order. While it is on, each card shows a drag
    // handle, the cards' own links and controls are inert (CSS) and their
    // clicks are swallowed (onClick), so a card can be grabbed without
    // following it.
    Alpine.data('cardDeckReorder', () => ({
      reordering: false,
      order: [],

      init() {
        this.order = cardDeckOrderIds(cardDeckApplyOrder(this.$root, true))
      },
      toggleReorder() {
        this.reordering = !this.reordering
        cardDeckSetControlsDisabled(this.$root, this.reordering)
      },
      // Renumbers the cards and announces the result. The hidden input, when
      // the deck has a Name, is bound to the same array.
      commit() {
        this.order = cardDeckOrderIds(cardDeckApplyOrder(this.$root, false))
        this.$root.dispatchEvent(new CustomEvent('popui-card-deck-reorder', {
          bubbles: true,
          detail: { order: this.order },
        }))
      },
      onPointerDown(e) {
        if (!this.reordering || e.button > 0) return
        const card = cardDeckOwnCard(this.$root, e.target)
        if (!card) return
        // Suppresses the text selection and native image drag that would
        // otherwise start under the pointer. It also keeps focus off the card,
        // which is deliberate: cards are dragged, never selected.
        e.preventDefault()
        cardDeckStartDrag(this.$root, card, e, () => this.commit())
      },
      // Escape is a way out of the mode, not a way to reorder — the cards
      // themselves are drag-only.
      onKeydown(e) {
        if (this.reordering && e.key === 'Escape') this.toggleReorder()
      },
      // Cards are grab handles while reordering, never links or buttons.
      onClick(e) {
        if (!this.reordering) return
        if (!cardDeckOwnCard(this.$root, e.target)) return
        e.preventDefault()
        e.stopPropagation()
      },
    }))
  })
})();
