import { networkSelector } from 'store/selectors/networkSelectors';
import { getState } from 'store/store';

export function getRewaLabel(state = getState()) {
  return networkSelector(state).rewaLabel;
}
