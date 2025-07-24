import { ExtensionProvider } from '@terradharitri/sdk-extension-provider';

export async function getExtensionProvider(address: string) {
  try {
    const provider = ExtensionProvider.getInstance().setAddress(address);

    const success = await provider.init();

    if (success) {
      return provider;
    } else {
      console.error(
        'Could not initialise extension, make sure Dharitri wallet extension is installed.'
      );
    }
  } catch (err) {
    console.error('Unable to login to ExtensionProvider', err);
  }
  return null;
}
