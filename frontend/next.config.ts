import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Hide the on-screen dev indicator widget.
  devIndicators: false,

  images: {
    // Item media is seeded with picsum.photos placeholders; /seed/ URLs
    // redirect to the fastly.picsum.photos CDN, so both hosts are allowed.
    remotePatterns: [
      { protocol: "https", hostname: "picsum.photos" },
      { protocol: "https", hostname: "fastly.picsum.photos" },
    ],
  },
};

export default nextConfig;
