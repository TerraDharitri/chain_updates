import { networkSelector } from 'store/selectors';
import { getState } from 'store/store';

export function getRewaLabel(state = getState()) {
  return networkSelector(state).rewaLabel;
}
