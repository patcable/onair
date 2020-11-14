source = ["./onair"]
bundle_id = "net.pcable.onair"

apple_id {
  username = "@env:AC_EMAIL"
  password = "@keychain:apple_id"
}

sign {
  application_identity = "B6AE7396AD644B78D62F4B970E92E661A7D97B44"
}

zip {
  output_path = "onair.zip"
}
