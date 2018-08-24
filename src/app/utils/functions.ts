import {Subscription} from 'rxjs';

export function clearSubs(subs$: Subscription[]) {
  subs$.forEach((sub: Subscription) => sub.unsubscribe());
}
