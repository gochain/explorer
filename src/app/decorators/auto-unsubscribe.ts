import {unsubArr} from '../utils/functions';

export function AutoUnsubscribe(target: string) {

  return function (constructor) {
    const original = constructor.prototype.ngOnDestroy;

    constructor.prototype.ngOnDestroy = function () {
      unsubArr(this[target]);
      if (original && typeof original === 'function') {
        original.apply(this, arguments);
      }
    };
  };
}

