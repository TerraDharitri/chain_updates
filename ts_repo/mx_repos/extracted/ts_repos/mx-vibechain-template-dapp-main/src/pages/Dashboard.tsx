import React, { useEffect } from "react";
import { getAccountProvider, useGetAccount } from "../lib";
import { useNavigate, useSearchParams } from "react-router-dom";

const Dashboard = () => {
  const navigate = useNavigate();
  const provider = getAccountProvider();
  const [searchParams] = useSearchParams();
  const accessToken = searchParams.get("access_token");

  useEffect(() => {
    if (accessToken) {
      console.log("Access Token:", accessToken);
    }
  }, [accessToken]);

  const handleLogout = async () => {
    await provider.logout();
    navigate("/unlock");
  };

  const { address } = useGetAccount();

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-2xl shadow-lg p-8">
            <div className="flex justify-between items-center mb-8">
              <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>
              <button
                onClick={handleLogout}
                className="px-6 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors duration-200"
              >
                Logout
              </button>
            </div>

            <div className="space-y-6">
              <div className="bg-gray-50 rounded-xl p-6">
                <h2 className="text-sm font-semibold text-gray-500 mb-2">
                  WALLET ADDRESS
                </h2>
                <p className="font-mono text-gray-800 break-all">{address}</p>
              </div>

              {accessToken && (
                <div className="bg-gray-50 rounded-xl p-6">
                  <h2 className="text-sm font-semibold text-gray-500 mb-2">
                    ACCESS TOKEN
                  </h2>
                  <p className="font-mono text-gray-800 break-all">
                    {accessToken}
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
