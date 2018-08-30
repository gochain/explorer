import {Subscription} from 'rxjs';

/**
 * clears array from subscriptions
 * @param {Array<Subscription>} arr
 */
export function unsubArr(arr: Array<Subscription>) {
  arr.forEach(item => item.unsubscribe());
}
