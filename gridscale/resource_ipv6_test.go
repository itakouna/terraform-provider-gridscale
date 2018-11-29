package gridscale

import (
	"fmt"
	"testing"

	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGridscaleIpv6_Basic(t *testing.T) {
	var object gsclient.Ip
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleIpv6Config_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleIpv6Exists("gridscale_ipv6.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_ipv6.foo", "name", name),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleIpv6Exists(n string, object *gsclient.Ip) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}

		client := testAccProvider.Meta().(*gsclient.Client)

		id := rs.Primary.ID

		foundObject, err := client.GetIp(id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUuid != id {
			return fmt.Errorf("Object not found")
		}

		*object = *foundObject

		return nil
	}
}

func testAccCheckDataSourceGridscaleIpv6Config_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_ipv6" "foo" {
  name   = "%s"
}
`, name)
}