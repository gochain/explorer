import {Subscription} from 'rxjs';

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
