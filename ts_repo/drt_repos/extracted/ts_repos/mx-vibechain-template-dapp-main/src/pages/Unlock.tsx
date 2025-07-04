import React from "react";
import {
  useGetLoginInfo,
  useGetAccount,
  ProviderFactory,
  ProviderTypeEnum,
} from "../lib";
import { useNavigate } from "react-router-dom";
import { IProviderFactory } from "../lib";

const Unlock = () => {
  const { isLoggedIn } = useGetLoginInfo();
  const { address } = useGetAccount();
  const navigate = useNavigate();

  React.useEffect(() => {
    if (isLoggedIn && address) {
      navigate("/dashboard");
    }
  }, [isLoggedIn, address, navigate]);

  const handleLogin = async ({ type, anchor }: IProviderFactory) => {
    const provider = await ProviderFactory.create({
      type,
      anchor,
    });
    // provider.getProvider().init();
    const loginData = await provider?.login();

    if (!loginData) return;

    navigate("/dashboard", {
      state: { from: "/unlock" },
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50 flex items-center justify-center">
      <div className="max-w-md w-full mx-4">
        <div className="bg-white rounded-2xl shadow-lg p-8 text-center">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-800 mb-2">
              Welcome to MultiversX
            </h1>
            <p className="text-gray-600">Connect your wallet to continue</p>
          </div>

          <div className="space-y-4">
            <button
              onClick={() =>
                handleLogin({ type: ProviderTypeEnum.crossWindow })
              }
              className="w-full py-3 px-4 bg-gray-600 hover:bg-gray-700 text-white rounded-lg font-medium transition-colors duration-200 flex items-center justify-center space-x-2"
            >
              <img
                src="https://wallet.multiversx.com/favicon.ico"
                alt="Web Wallet"
                className="w-5 h-5"
              />
              <span>Connect with Web Wallet</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Unlock;
