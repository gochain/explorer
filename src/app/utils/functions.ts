import {Subscription} from 'rxjs';

declare const window: any;

/**
 * clears array from subscriptions
 * @param {Array<Subscription>} arr
 */
export function unsubArr(arr: Array<Subscription>) {
  arr.forEach(item => item.unsubscribe());
}

/**
 * recursively checks object keys
 * @param obj
 * @param keys
 */
export function objHas(obj: any, keys: string): boolean {
  return !!keys.split('.').reduce((acc, key) => acc.hasOwnProperty(key) ? (acc[key] || 1) : false, obj);
}

/**
 *
 */
export function isPrivateMode(): boolean {
  const testLocalStorage = () => {
    try {
      if (localStorage.length) {
        return true;
      } else {
        localStorage.x = 1;
        localStorage.removeItem('x');
        return true;
      }
    } catch (e) {
      // Safari only enables cookie in private mode
      // if cookie is disabled then all client side storage is disabled
      // if all client side storage is disabled, then there is no point
      // in using private mode
      return !navigator.cookieEnabled;
    }
  };
  // Chrome & Opera
  if (window.webkitRequestFileSystem) {
    return window.webkitRequestFileSystem(0, 0, true, false);
  }
  // Firefox
  if ('MozAppearance' in document.documentElement.style) {
    const db = indexedDB.open('test');
    db.onerror = () => {
      return true;
    };
    db.onsuccess = () => {
      return false;
    };
  }
  // Safari
  if (/constructor/i.test(window.HTMLElement)) {
    return testLocalStorage();
  }
  // IE10+ & Edge
  if (!window.indexedDB && (window.PointerEvent || window.MSPointerEvent)) {
    return true;
  }
  // others
  return false;
}
