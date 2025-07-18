import { rewaLabelSelector } from 'reduxStore/selectors';
import { store } from 'reduxStore/store';

export function getRewaLabel() {
  return rewaLabelSelector(store.getState());
}
