class Cix < Formula
  desc "Run and debug GitLab CI pipelines locally"
  homepage "https://github.com/tiep-cixplatform/cix"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/tiep-cixplatform/cix/releases/download/v1.0.0/cix-darwin-arm64"
      sha256 "2a0ff9f7d22ce81886743f3de9e10820ba8ecd7f549751ca1cade57be54e10b0"
    else
      url "https://github.com/tiep-cixplatform/cix/releases/download/v1.0.0/cix-darwin-amd64"
      sha256 "ab44cf6fa087e5299c75d290eb178545cf3fffb7e6921b3b8feb0c9e89cce8df"
    end
  end

  def install
    if Hardware::CPU.arm?
      bin.install "cix-darwin-arm64" => "cix"
    else
      bin.install "cix-darwin-amd64" => "cix"
    end
  end

  test do
    system "#{bin}/cix", "--version"
  end
end
